package cli

import (
	"log"

	"github.com/hansung080/gchain/node"
)

func createBlockchain(nodeID, addr string) {
	if !node.ValidateAddress(addr) {
		log.Panicf("Invalid address: %v\n", addr)
	}

	bc := node.CreateBlockchain(nodeID, addr)
	defer bc.Close()

	node.UTXOSet{bc}.Reindex()
}
