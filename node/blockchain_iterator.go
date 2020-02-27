package node

import (
	"log"
	"github.com/boltdb/bolt"
)

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (i *BlockchainIterator) Next() *Block {
	var block *Block

	if err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = UnmarshalBlock(encodedBlock)
		return nil

	}); err != nil {
		log.Panic(err)
	}

	i.currentHash = block.PrevHash
	return block
}

func (i *BlockchainIterator) HasNext() bool {
	if len(i.currentHash) == 0 {
		return false
	}
	return true
}
