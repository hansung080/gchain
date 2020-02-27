package node

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"
)

/**
  @ Transaction Input/Output Diagram
    - A input must be connected with a output of the previous transaction.
    - Outputs which are not connected with inputs could exist. Those outputs are called UTXO.
    - Inputs of a transaction could be connected with outputs of many transactions.

         tx0
    -------------
    | in0  out0 | ---                   |
    | in1  out1 |   |                   |
    -------------   |                   |
         tx1        |        tx3        |        tx4
    -------------   |   -------------   |   -------------
    | in0  out0 |   --> | in0  out0 |   --> | in0  out0 |
    |      out1 | ----> | in1  out1 | ----> | in1       |
    -------------   --> | in2       |       -------------
         tx2        |   -------------
    -------------   |
    | in0       |   |
    | in1  out0 | ---
    | in2       |
    -------------
*/

const subsidy = 10

type Transaction struct {
	ID    []byte  // transaction ID
	Vins  []TxIn  // transaction input list
	Vouts []TxOut // transaction output list
}

func (tx Transaction) Marshal() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	if err := encoder.Encode(tx); err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

func (tx *Transaction) Hash() []byte {
	copiedTx := *tx
	copiedTx.ID = []byte{}
	hash := sha256.Sum256(copiedTx.Marshal())
	return hash[:]
}

func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vins) == 1 && len(tx.Vins[0].Txid) == 0 && tx.Vins[0].Vout == -1
}

func (tx *Transaction) TrimmedCopy() Transaction {
	var ins []TxIn
	var outs []TxOut

	for _, in := range tx.Vins {
		ins = append(ins, TxIn{
			Txid: in.Txid,
			Vout: in.Vout,
			Sig:  nil,
			Pkey: nil,
		})
	}

	for _, out := range tx.Vouts {
		outs = append(outs, TxOut{
			Value:    out.Value,
			PkeyHash: out.PkeyHash,
		})
	}

	return Transaction{
		ID:    tx.ID,
		Vins:  ins,
		Vouts: outs,
	}
}

func (tx *Transaction) Sign(prevTxs map[string]Transaction, skey ecdsa.PrivateKey) {
	if tx.IsCoinbase() {
		return
	}

	for _, in := range tx.Vins {
		if prevTxs[hex.EncodeToString(in.Txid)].ID == nil {
			log.Panic("Invalid previous transaction ID")
		}
	}

	copiedTx := tx.TrimmedCopy()
	for idx, in := range copiedTx.Vins {
		prevTx := prevTxs[hex.EncodeToString(in.Txid)]
		copiedTx.Vins[idx].Sig = nil
		copiedTx.Vins[idx].Pkey = prevTx.Vouts[in.Vout].PkeyHash

		data := fmt.Sprintf("%x\n", copiedTx)

		// TODO: check that the sign data which is not a hash would be a problem.
		r, s, err := ecdsa.Sign(rand.Reader, &skey, []byte(data))
		if err != nil {
			log.Panic(err)
		}
		sig := append(r.Bytes(), s.Bytes()...)

		tx.Vins[idx].Sig = sig
		copiedTx.Vins[idx].Pkey = nil
	}
}

func (tx *Transaction) Verify(prevTxs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	for _, in := range tx.Vins {
		if prevTxs[hex.EncodeToString(in.Txid)].ID == nil {
			log.Panic("Invalid previous transaction ID")
		}
	}

	copiedTx := tx.TrimmedCopy()
	curve := elliptic.P256()

	for idx, in := range tx.Vins {
		prevTx := prevTxs[hex.EncodeToString(in.Txid)]
		copiedTx.Vins[idx].Sig = nil
		copiedTx.Vins[idx].Pkey = prevTx.Vouts[in.Vout].PkeyHash

		r := big.Int{}
		s := big.Int{}
		sigHalfLen := len(in.Sig) / 2
		r.SetBytes(in.Sig[:sigHalfLen])
		s.SetBytes(in.Sig[sigHalfLen:])

		x := big.Int{}
		y := big.Int{}
		pkeyHalfLen := len(in.Pkey) / 2
		x.SetBytes(in.Pkey[:pkeyHalfLen])
		y.SetBytes(in.Pkey[pkeyHalfLen:])

		data := fmt.Sprintf("%x\n", copiedTx)
		pkey := ecdsa.PublicKey{
			Curve: curve,
			X:     &x,
			Y:     &y,
		}

		if !ecdsa.Verify(&pkey, []byte(data), &r, &s) {
			return false
		}

		copiedTx.Vins[idx].Pkey = nil
	}

	return true
}

func (tx Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf(" - transaction %x", tx.ID))

	for i, in := range tx.Vins {
		lines = append(lines, fmt.Sprintf("     input %d", i))
		lines = append(lines, fmt.Sprintf("       txid: %x", in.Txid))
		lines = append(lines, fmt.Sprintf("       out: %d", in.Vout))
		lines = append(lines, fmt.Sprintf("       sig: %x", in.Sig))
		lines = append(lines, fmt.Sprintf("       pkey: %x", in.Pkey))
	}

	for i, out := range tx.Vouts {
		lines = append(lines, fmt.Sprintf("     output %d", i))
		lines = append(lines, fmt.Sprintf("       value: %d", out.Value))
		lines = append(lines, fmt.Sprintf("       pkeyHash: %x", out.PkeyHash))
	}

	return strings.Join(lines, "\n")
}

func NewTransaction(wallet *Wallet, to string, amount int, utxoSet *UTXOSet) *Transaction {
	var inputs []TxIn
	var outputs []TxOut

	sum, utxos := utxoSet.FindSpendableOuts(HashPkey(wallet.Pkey), amount)
	if sum < amount {
		log.Panic("Balance not enough")
	}

	for txid, outs := range utxos {
		txidBytes, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		// make a input list.
		// make inputs to spend UTXOs.
		for _, out := range outs {
			inputs = append(inputs, TxIn{
				Txid: txidBytes,
				Vout: out,
				Sig:  nil,
				Pkey: wallet.Pkey,
			})
		}
	}

	// make a output list.
	// make a output to give coins.
	outputs = append(outputs, *NewTxOut(amount, to))

	// make a output to get the change back, because a output is indivisible.
	if sum > amount {
		from := string(wallet.GetAddress())
		outputs = append(outputs, *NewTxOut(sum - amount, from))
	}

	// make a transaction.
	tx := Transaction{
		ID:    nil,
		Vins:  inputs,
		Vouts: outputs,
	}

	tx.ID = tx.Hash()
	utxoSet.BC.SignTx(&tx, wallet.Skey)
	return &tx
}

func NewCoinbaseTx(to, data string) *Transaction {
	if data == "" {
		randData := make([]byte, 20)
		if _, err := rand.Read(randData); err != nil {
			log.Panic(err)
		}

		data = fmt.Sprintf("%x", randData)
	}

	in := TxIn{
		Txid: []byte{},
		Vout: -1,
		Sig:  nil,
		Pkey: []byte(data),
	}

	out := *NewTxOut(subsidy, to)
	
	tx := Transaction{
		ID:    nil,
		Vins:  []TxIn{in},
		Vouts: []TxOut{out},
	}

	tx.ID = tx.Hash()
	return &tx
}

func UnmarshalTx(data []byte) Transaction {
	var tx Transaction

	decoder := gob.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&tx); err != nil {
		log.Panic(err)
	}

	return tx
}
