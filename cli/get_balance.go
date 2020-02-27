package cli

import (
	"fmt"
	"log"

	"github.com/hansung080/gchain/node"
)

func getBalance(nodeID, addr string) {
	if !node.ValidateAddress(addr) {
		log.Panicf("Invalid address: %v\n", addr)
	}

	bc := node.NewBlockchain(nodeID)
	defer bc.Close()

	pkeyHash := node.GetPkeyHashFromAddress([]byte(addr))
	utxos := node.UTXOSet{bc}.FindUTXOs(pkeyHash)

	balance := 0
	for _, out := range utxos {
		balance += out.Value
	}

	fmt.Println(balance)
}
