package server

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"github.com/hansung080/gchain/node"
)

func handleConnection(conn net.Conn, bc *node.Blockchain) {
	defer conn.Close()

	req, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}

	cmd := bytesToCommand(req[:commandLen])
	fmt.Printf("Received command: %s\n", cmd)

	switch cmd {
	case "addr":
		handleAddress(req)
	case "block":
		handleBlock(req, bc)
	case "getblocks":
		handleGetBlocks(req, bc)
	case "getdata":
		handleGetData(req, bc)
	case "inv":
		handleInventory(req, bc)
	case "tx":
		handleTx(req, bc)
	case "version":
		handleVersion(req, bc)
	default:
		fmt.Printf("Invalid command: %s\n", cmd)
	}
}

func handleAddress(req []byte) {
	var payload address

	unmarshalGob(req[commandLen:], &payload)
	knownAddrs = append(knownAddrs, payload.Addrs...)
	fmt.Printf("Known addresses count: %d", len(knownAddrs))
	requestBlocks()
}

func handleBlock(req []byte, bc *node.Blockchain) {
	var payload block

	unmarshalGob(req[commandLen:], &payload)
	block := node.UnmarshalBlock(payload.Block)
	fmt.Printf("Received a new block: %x", block.Hash)
	bc.AddBlock(block)

	if len(blocksInTransit) > 0 {
		blockHash := blocksInTransit[0]
		sendGetData(payload.From, "block", blockHash)
		blocksInTransit = blocksInTransit[1:]
	} else {
		node.UTXOSet{bc}.Reindex()
	}
}

func handleGetBlocks(req []byte, bc *node.Blockchain) {
	var payload getblocks

	unmarshalGob(req[commandLen:], &payload)
	blockHashes := bc.GetBlockHashes()
	sendInventory(payload.From, "block", blockHashes)
}

func handleGetData(req []byte, bc *node.Blockchain) {
	var payload getdata

	unmarshalGob(req[commandLen:], &payload)

	if payload.Type == "block" {
		block, err := bc.GetBlock(payload.ID)
		if err != nil {
			return
		}

		sendBlock(payload.From, &block)

	} else if payload.Type == "tx" {
		txid := hex.EncodeToString(payload.ID)
		tx := mempool[txid]
		sendTx(payload.From, &tx)
		//delete(mempool, txid)
	}
}

func handleInventory(req []byte, bc *node.Blockchain) {
	var payload inventory

	unmarshalGob(req[commandLen:], &payload)
	fmt.Printf("Received inventory: type: %s, items: %d\n", payload.Type, len(payload.Items))

	if payload.Type == "block" {
		blocksInTransit = payload.Items
		blockHash := payload.Items[0]
		sendGetData(payload.From, "block", blockHash)

		newInTransit := [][]byte{}
		for _, b := range blocksInTransit {
			if bytes.Compare(b, blockHash) != 0 {
				newInTransit = append(newInTransit, b)
			}
		}

		blocksInTransit = newInTransit

	} else if payload.Type == "tx" {
		txid := payload.Items[0]
		if _, exist := mempool[hex.EncodeToString(txid)]; !exist {
			sendGetData(payload.From, "tx", txid)
		}
	}
}

func handleTx(req []byte, bc *node.Blockchain) {
	var payload transaction

	unmarshalGob(req[commandLen:], &payload)
	tx := node.UnmarshalTx(payload.Tx)
	mempool[hex.EncodeToString(tx.ID)] = tx

	if nodeAddr == knownAddrs[0] {
		for _, addr := range knownAddrs {
			if addr != nodeAddr && addr != payload.From {
				sendInventory(addr, "tx", [][]byte{tx.ID})
			}
		}
	} else {
		if len(mempool) >= 2 && len(minerAddr) > 0 {
		MineTransactions:
			var txs []*node.Transaction
			for id := range mempool {
				tx := mempool[id]
				if bc.VerifyTx(&tx) {
					txs = append(txs, &tx)
				}
			}

			if len(txs) < 1 {
				fmt.Println("All transactions are failed to verify.")
				return
			}

			coinbase := node.NewCoinbaseTx(minerAddr, "")
			txs = append(txs, coinbase)

			newBlock := bc.MineBlock(txs)
			node.UTXOSet{bc}.Reindex()
			fmt.Println("Mined a new block.")

			for _, tx := range txs {
				delete(mempool, hex.EncodeToString(tx.ID))
			}

			for _, addr := range knownAddrs {
				if addr != nodeAddr {
					sendInventory(addr, "block", [][]byte{newBlock.Hash})
				}
			}

			if len(mempool) > 0 {
				goto MineTransactions
			}
		}
	}
}

func handleVersion(req []byte, bc *node.Blockchain) {
	var payload version

	unmarshalGob(req[commandLen:], &payload)
	myHeight := bc.GetBestHeight()
	yourHeight := payload.BestHeight

	if myHeight < yourHeight {
		sendGetBlocks(payload.From)
	} else if myHeight > yourHeight {
		sendVersion(payload.From, bc)
	}

	//sendAddress(payload.From)
	if !isNodeKnown(payload.From) {
		knownAddrs = append(knownAddrs, payload.From)
	}
}
