package server

import (
	"log"
	"net"
	"fmt"
	"io"
	"bytes"

	"github.com/hansung080/gchain/node"
)

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

func sendBlock(addr string, b *node.Block) {
	payload := block{
		From:  nodeAddr,
		Block: b.Marshal(),
	}

	resp := append(commandToBytes("block"), marshalGob(payload)...)
	send(addr, resp)
}

func sendGetBlocks(addr string) {
	payload := getblocks{nodeAddr}
	resp := append(commandToBytes("getblocks"), marshalGob(payload)...)
	send(addr, resp)
}

func sendGetData(addr, typ string, id []byte) {
	payload := getdata{
		From: nodeAddr,
		Type: typ,
		ID:   id,
	}

	resp := append(commandToBytes("getdata"), marshalGob(payload)...)
	send(addr, resp)
}

func sendInventory(addr, typ string, items [][]byte) {
	payload := inventory{
		From:  nodeAddr,
		Type:  typ,
		Items: items,
	}

	resp := append(commandToBytes("inv"), marshalGob(payload)...)
	send(addr, resp)
}

func sendTx(addr string, tx *node.Transaction) {
	payload := transaction{
		From: nodeAddr,
		Tx:   tx.Marshal(),
	}

	resp := append(commandToBytes("tx"), marshalGob(payload)...)
	send(addr, resp)
}

func sendVersion(addr string, bc *node.Blockchain) {
	payload := version{
		From:       nodeAddr,
		Version:    nodeVersion,
		BestHeight: bc.GetBestHeight(),
	}

	resp := append(commandToBytes("version"), marshalGob(payload)...)
	send(addr, resp)
}

func send(addr string, resp []byte) {
	conn, err := net.Dial(protocol, addr)
	if err != nil {
		fmt.Printf("Cannot create connection: %s\n", addr)
		var updatedAddrs []string
		for _, a := range knownAddrs {
			if a != addr {
				updatedAddrs = append(updatedAddrs, a)
			}
		}

		knownAddrs = updatedAddrs
		return
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(resp))
	if err != nil {
		log.Panic(err)
	}
}
