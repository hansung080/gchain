package node

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	Timestamp int64
	Txs       []*Transaction
	PrevHash  []byte
	Hash      []byte
	Nonce     int
	Height    int
}

func (b *Block) Marshal() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	if err := encoder.Encode(b); err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

func (b *Block) HashTxs() []byte {
	var txs [][]byte

	for _, tx := range b.Txs {
		txs = append(txs, tx.Marshal())
	}

	tree := NewMerkleTree(txs)
	return tree.Root.Hash
}

func NewBlock(txs []*Transaction, prevHash []byte, height int) *Block {
	block := &Block{
		Timestamp: time.Now().Unix(),
		Txs:       txs,
		PrevHash:  prevHash,
		Hash:      []byte{},
		Nonce:     0,
		Height:    height,
	}

	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{}, 0)
}

func UnmarshalBlock(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&block); err != nil {
		log.Panic(err)
	}

	return &block
}
