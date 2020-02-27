package cli

import (
	"fmt"

	"github.com/hansung080/gchain/node"
)

func createWallet(nodeID string) {
	wallets, err := node.NewWallets(nodeID)
	if err != nil {
		fmt.Printf("Wallet creation failure: %s\n", err.Error())
		return
	}

	addr := wallets.CreateWallet()
	wallets.SaveFile(nodeID)
	fmt.Println(addr)
}
