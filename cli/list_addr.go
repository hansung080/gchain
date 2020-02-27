package cli

import (
	"log"
	"fmt"

	"github.com/hansung080/gchain/node"
)

func listAddresses(nodeID string) {
	wallets, err := node.NewWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}

	addrs := wallets.GetAddresses()
	for _, addr := range addrs {
		fmt.Println(addr)
	}
}
