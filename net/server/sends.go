package server

import "github.com/hansung080/gchain/node"

func requestBlocks() {
	for _, addr := range knownAddrs {
		sendGetBlocks(addr)
	}
}

func sendAddress(addr string) {
	payload := address{knownAddrs}
	payload.Addrs = append(payload.Addrs, nodeAddr)
	resp := append(commandToBytes("addr"), marshalGob(payload)...)
	send(addr, resp)
}

func sendBlock(addr string, block *node.Block) {
	// khs working here...
}

func sendGetBlocks(addr string) {

}

func sendGetData(addr, typ string, id []byte) {

}

func sendInventory(addr, typ string, items [][]byte) {

}

func sendTx(addr string, tx *node.Transaction) {

}

func sendVersion(addr string, bc *node.Blockchain) {

}

func send(addr string, resp []byte) {

}
