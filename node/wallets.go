package node

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
)

const walletFile = "wallet_%s.dat"

type Wallets struct {
	Wallets map[string]*Wallet
}

func (ws *Wallets) CreateWallet() string {
	wallet := NewWallet()
	addr := string(wallet.GetAddress())
	ws.Wallets[addr] = wallet
	return addr
}

func (ws Wallets) GetWallet(addr string) Wallet {
	return *ws.Wallets[addr]
}

func (ws *Wallets) GetAddresses() []string {
	var addrs []string

	for addr := range ws.Wallets {
		addrs = append(addrs, addr)
	}

	return addrs
}

func (ws Wallets) SaveFile(nodeID string) {
	walletFile := fmt.Sprintf(walletFile, nodeID)

	var content bytes.Buffer
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	if err := encoder.Encode(ws); err != nil {
		log.Panic(err)
	}

	if err := ioutil.WriteFile(walletFile, content.Bytes(), 0644); err != nil {
		log.Panic(err)
	}
}

func (ws *Wallets) LoadFile(nodeID string) error {
	walletFile := fmt.Sprintf(walletFile, nodeID)
	if !FileExist(walletFile) {
		return nil // No error return here, because wallet file does not exist when the first wallet is created.
	}

	content, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}

	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(content))
	if err := decoder.Decode(&wallets); err != nil {
		log.Panic(err)
	}

	ws.Wallets = wallets.Wallets
	return nil
}

func NewWallets(nodeID string) (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)
	err := wallets.LoadFile(nodeID)
	return &wallets, err
}
