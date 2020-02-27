package cli

import (
	"log"

	"github.com/hansung080/gchain/node"
)

func send(nodeID, from, to string, amount int, mine bool) {
	if !node.ValidateAddress(from) {
		log.Panicf("Invalid address: %v\n", from)
	}

	if !node.ValidateAddress(to) {
		log.Panicf("Invalid address: %v\n", to)
	}

	bc := node.NewBlockchain(nodeID)
	defer bc.Close()
	utxoSet := node.UTXOSet{bc}

	wallets, err := node.NewWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}
	wallet := wallets.GetWallet(from)

	tx := node.NewTransaction(&wallet, to, amount, &utxoSet)
	if mine {
		coinbase := node.NewCoinbaseTx(from, "")
		txs := []*node.Transaction{coinbase, tx}
		block := bc.MineBlock(txs)
		utxoSet.Update(block)
	} else {
		// TODO: send transaction to another node.
	}
}