package node

import (
	"log"
	"encoding/hex"

	"github.com/boltdb/bolt"
)

const utxoBucket = "chainstate"

type UTXOSet struct {
	BC *Blockchain
}

func (u UTXOSet) FindUTXOs(pkeyHash []byte) []TxOut {
	var utxos []TxOut

	if err := u.BC.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			outs := UnmarshalOuts(v)
			for _, out := range outs.Outs {
				if out.LockedWith(pkeyHash) {
					utxos = append(utxos, out)
				}
			}
		}

		return nil

	}); err != nil {
		log.Panic(err)
	}

	return utxos
}

func (u UTXOSet) FindSpendableOuts(pkeyHash []byte, amount int) (int, map[string][]int) {
	sum := 0
	utxos := make(map[string][]int)

	if err := u.BC.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			txid := hex.EncodeToString(k)
			outs := UnmarshalOuts(v)
			for idx, out := range outs.Outs {
				if out.LockedWith(pkeyHash) && sum < amount {
					sum += out.Value
					utxos[txid] = append(utxos[txid], idx)
				}
			}
		}

		return nil

	}); err != nil {
		log.Panic(err)
	}
	
	return sum, utxos
}

func (u UTXOSet) CountTxs() int {
	count := 0

	if err := u.BC.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			count++
		}

		return nil

	}); err != nil {
		log.Panic(err)
	}

	return count
}

func (u UTXOSet) Reindex() {
	bucketName := []byte(utxoBucket)

	if err := u.BC.db.Update(func(tx *bolt.Tx) error {
		if err := tx.DeleteBucket(bucketName); err != nil && err != bolt.ErrBucketNotFound {
			log.Panic(err)
		}

		if _, err := tx.CreateBucket(bucketName); err != nil {
			log.Panic(err)
		}

		return nil

	}); err != nil {
		log.Panic(err)
	}

	utxos := u.BC.FindUTXOs()

	if err := u.BC.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		for txid, outs := range utxos {
			key, err := hex.DecodeString(txid)
			if err != nil {
				log.Panic(err)
			}

			if err = b.Put(key, outs.Marshal()); err != nil {
				log.Panic(err)
			}
		}

		return nil

	}); err != nil {
		log.Panic(err)
	}
}

func (u UTXOSet) Update(block *Block) {
	if err := u.BC.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))

		for _, tx := range block.Txs {
			// update rebuilt UTXOs of previous transactions
			if !tx.IsCoinbase() {
				for _, in := range tx.Vins {
					newOuts := TxOuts{}
					oldOuts := UnmarshalOuts(b.Get(in.Txid))
					for idx, out := range oldOuts.Outs {
						if idx != in.Vout {
							newOuts.Outs = append(newOuts.Outs, out)
						}
					}

					if len(newOuts.Outs) == 0 {
						if err := b.Delete(in.Txid); err != nil {
							log.Panic(err)
						}
					} else {
						if err := b.Put(in.Txid, newOuts.Marshal()); err != nil {
							log.Panic(err)
						}
					}
				}
			}

			// update new UTXOs of current transaction
			newOuts := TxOuts{}
			for _, out := range tx.Vouts {
				newOuts.Outs = append(newOuts.Outs, out)
			}

			if err := b.Put(tx.ID, newOuts.Marshal()); err != nil {
				log.Panic(err)
			}
		}

		return nil

	}); err != nil {
		log.Panic(err)
	}
}
