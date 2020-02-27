package node

import (
	"fmt"
	"math/big"
	"bytes"
	"math"
	"crypto/sha256"
	"time"
)

const (
	targetBits = 16 // 24 // targetBits gets larger, target gets smaller, and the difficulty of POW gets higher.
	maxNonce   = math.MaxInt64
)

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	return bytes.Join([][]byte{
		pow.block.PrevHash,
		pow.block.HashTxs(),
		IntToBytes(pow.block.Timestamp),
		IntToBytes(int64(targetBits)),
		IntToBytes(int64(nonce)),
	}, []byte{})
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	//fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)
	fmt.Println("Mining the block...")
	startTime := time.Now()
	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target) == -1 {
			break
		}
		nonce++
	}

	elapsedTime := time.Since(startTime)
	fmt.Println()
	fmt.Println("Done: nonce:", nonce, ", elapsed time:", elapsedTime)
	fmt.Println()
	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	return hashInt.Cmp(pow.target) == -1
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256 - targetBits))
	return &ProofOfWork{
		block:  b,
		target: target,
	} 
}
