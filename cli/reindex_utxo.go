package cli

import (
	"fmt"

	"github.com/hansung080/gchain/node"
)

func reindexUTXO(nodeID string) {
	bc := node.NewBlockchain(nodeID)
	defer bc.Close()

	utxoSet := node.UTXOSet{bc}
	utxoSet.Reindex()
	count := utxoSet.CountTxs()
	fmt.Printf("%d txs in UTXO set\n", count)
}
