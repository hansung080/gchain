package node

import (
	"fmt"
	"log"
	"encoding/hex"
	"os"
	"crypto/ecdsa"
	"bytes"
	"errors"

	"github.com/boltdb/bolt"
)

const (
	dbFile       = "blockchain_%s.db"
	blocksBucket = "blocks"
	genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"
)

type Blockchain struct {
	tip []byte // last block hash
	db  *bolt.DB
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{
		currentHash: bc.tip,
		db:          bc.db,
	}
}

func (bc *Blockchain) FindUTXOs() map[string]TxOuts {
	utxos := make(map[string]TxOuts)
	stxos := make(map[string][]int)
	iter := bc.Iterator()

	for iter.HasNext() {
		block := iter.Next()
		for _, tx := range block.Txs {
			txid := hex.EncodeToString(tx.ID)

		Outputs:
			for idx, out := range tx.Vouts {
				for _, stxo := range stxos[txid] {
					if stxo == idx {
						continue Outputs
					}
				}

				outs := utxos[txid]
				outs.Outs = append(outs.Outs, out)
				utxos[txid] = outs
			}

			if !tx.IsCoinbase() {
				for _, in := range tx.Vins {
					inTxid := hex.EncodeToString(in.Txid)
					stxos[inTxid] = append(stxos[inTxid], in.Vout)
				}
			}
		}
	}

	return utxos
}

func (bc *Blockchain) FindTx(id []byte) (Transaction, error) {
	iter := bc.Iterator()
	for iter.HasNext() {
		block := iter.Next()
		for _, tx := range block.Txs {
			if bytes.Compare(tx.ID, id) == 0 {
				return *tx, nil
			}
		}
	}

	return Transaction{}, errors.New("Transaction not found")
}

func (bc *Blockchain) SignTx(tx *Transaction, skey ecdsa.PrivateKey) {
	prevTxs := make(map[string]Transaction)

	for _, in := range tx.Vins {
		prevTx, err := bc.FindTx(in.Txid)
		if err != nil {
			log.Panic(err)
		}

		prevTxs[hex.EncodeToString(prevTx.ID)] = prevTx
	}

	tx.Sign(prevTxs, skey)
}

func (bc *Blockchain) VerifyTx(tx *Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	prevTxs := make(map[string]Transaction)

	for _, in := range tx.Vins {
		prevTx, err := bc.FindTx(in.Txid)
		if err != nil {
			log.Panic(err)
		}

		prevTxs[hex.EncodeToString(prevTx.ID)] = prevTx
	}

	return tx.Verify(prevTxs)
}

func (bc *Blockchain) MineBlock(txs []*Transaction) *Block {
	for _, tx := range txs {
		if !bc.VerifyTx(tx) {
			log.Panic("Transaction verification failure")
		}
	}

	var lastHash []byte
	var lastHeight int
	if err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))
		lastBlock := UnmarshalBlock(b.Get(lastHash))
		lastHeight = lastBlock.Height
		return nil

	}); err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(txs, lastHash, lastHeight + 1)

	if err := bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if err := b.Put(newBlock.Hash, newBlock.Marshal()); err != nil {
			log.Panic(err)
		}

		if err := b.Put([]byte("l"), newBlock.Hash); err != nil {
			log.Panic(err)
		}

		bc.tip = newBlock.Hash
		return nil

	}); err != nil {
		log.Panic(err)
	}

	return newBlock
}

func (bc *Blockchain) AddBlock(block *Block) {
	if err := bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if b.Get(block.Hash) != nil {
			fmt.Println("Block already exists.")
			return nil
		}

		if err := b.Put(block.Hash, block.Marshal()); err != nil {
			log.Panic(err)
		}

		lastHash := b.Get([]byte("l"))
		lastBlock := UnmarshalBlock(b.Get(lastHash))
		if block.Height > lastBlock.Height {
			if err := b.Put([]byte("l"), block.Hash); err != nil {
				log.Panic(err)
			}
			bc.tip = block.Hash
		}

		return nil

	}); err != nil {
		log.Panic(err)
	}
}

func (bc *Blockchain) GetBlock(hash []byte) (Block, error) {
	var block Block

	if err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		blockBytes := b.Get(hash)
		if blockBytes == nil {
			return errors.New("Block not found")
		}

		block = *UnmarshalBlock(blockBytes)
		return nil

	}); err != nil {
		return block, err
	}

	return block, nil
}

func (bc *Blockchain) GetBlockHashes() [][]byte {
	var hashes [][]byte

	iter := bc.Iterator()
	for iter.HasNext() {
		block := iter.Next()
		hashes = append(hashes, block.Hash)
	}

	return hashes
}

func (bc *Blockchain) GetBestHeight() int {
	var lastHeight int

	if err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash := b.Get([]byte("l"))
		lastBlock := UnmarshalBlock(b.Get(lastHash))
		lastHeight = lastBlock.Height
		return nil

	}); err != nil {
		log.Panic(err)
	}

	return lastHeight
}

func (bc *Blockchain) Close() {
	bc.db.Close()
}

func CreateBlockchain(nodeID, addr string) *Blockchain {
	dbFile := fmt.Sprintf(dbFile, nodeID)
	if FileExist(dbFile) {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	coinbase := NewCoinbaseTx(addr, genesisCoinbaseData)
	genesis := NewGenesisBlock(coinbase)

	var tip []byte
	if err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic(err)
		}

		if err = b.Put(genesis.Hash, genesis.Marshal()); err != nil {
			log.Panic(err)
		}

		if err = b.Put([]byte("l"), genesis.Hash); err != nil {
			log.Panic(err)
		}

		tip = genesis.Hash
		return nil

	}); err != nil {
		log.Panic(err)
	}

	return &Blockchain{
		tip: tip,
		db:  db,
	}
}

func NewBlockchain(nodeID string) *Blockchain {
	dbFile := fmt.Sprintf(dbFile, nodeID)
	if !FileExist(dbFile) {
		fmt.Println("Blockchain not found. Create one first.")
		os.Exit(1)
	}

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	var tip []byte
	if err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte("l"))
		return nil

	}); err != nil {
		log.Panic(err)
	}

	return &Blockchain{
		tip: tip,
		db:  db,
	}
}
