package node

import (
	"bytes"
	"encoding/gob"
	"log"

	"github.com/hansung080/gchain/encoding/base58"
)

type TxOut struct {
	Value    int    // output value. generally coin
	PkeyHash []byte // output owner. lock output with this public key hash extracted from address user inputs
}

func (out *TxOut) Lock(addr []byte) {
	payload := base58.Decode(addr)
	out.PkeyHash = payload[addressVersionLen:len(payload) - addressChecksumLen]
}

func (out *TxOut) LockedWith(pkeyHash []byte) bool {
	return bytes.Compare(out.PkeyHash, pkeyHash) == 0
}

func NewTxOut(value int, addr string) *TxOut {
	txo := &TxOut{
		Value:    value,
		PkeyHash: nil,
	}

	txo.Lock([]byte(addr))
	return txo
}

type TxOuts struct {
	Outs []TxOut
}

func (outs TxOuts) Marshal() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	if err := encoder.Encode(outs); err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

func UnmarshalOuts(data []byte) TxOuts {
	var outs TxOuts

	decoder := gob.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&outs); err != nil {
		log.Panic(err)
	}

	return outs
}
