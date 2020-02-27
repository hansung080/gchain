package cli

import (
	"fmt"
	"strconv"

	"github.com/hansung080/gchain/node"
)

func printChain(nodeID string) {
	bc := node.NewBlockchain(nodeID)
	defer bc.Close()

	iter := bc.Iterator()
	for iter.HasNext() {
		block := iter.Next()

		fmt.Printf(" @ Block %d\n", block.Height)
		fmt.Printf(" - prev. hash: %x\n", block.PrevHash)
		fmt.Printf(" - hash: %x\n", block.Hash)
		fmt.Printf(" - pow: %s\n", strconv.FormatBool(node.NewProofOfWork(block).Validate()))
		for _, tx := range block.Txs {
			fmt.Println(tx)
		}
		fmt.Println()
	}
}
