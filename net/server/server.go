package server

import (
	"fmt"
	"net"
	"log"

	"github.com/hansung080/gchain/node"
)

const (
	protocol    = "tcp"
	nodeVersion = 1
	commandLen  = 12
)

var (
	nodeAddr  string
	minerAddr string
	knownAddrs      = []string{"localhost:3000"}
	blocksInTransit = [][]byte{}
	mempool         = make(map[string]node.Transaction)
)

type address struct {
	Addrs []string
}

type block struct {
	From  string
	Block []byte
}

type getblocks struct {
	From string
}

type getdata struct {
	From string
	Type string
	ID   []byte
}

type inventory struct {
	From  string
	Type  string
	Items [][]byte
}

type transaction struct {
	From string
	Tx   []byte
}

type version struct {
	From       string
	Version    int
	BestHeight int
}

func Start(nodeID, miner string) {
	nodeAddr = fmt.Sprintf("localhost:%s", nodeID)
	minerAddr = miner

	ln, err := net.Listen(protocol, nodeAddr)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	bc := node.NewBlockchain(nodeID)

	if nodeAddr != knownAddrs[0] {
		sendVersion(knownAddrs[0], bc)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}

		go handleConnection(conn, bc)
	}
}
