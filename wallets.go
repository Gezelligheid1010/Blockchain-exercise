package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
)

const walletFile = "wallet_%s.dat"

type Wallets struct {
	Wallets map[string]*Wallet
}

// SerializableWallet 解决gob: type elliptic.p256Curve has no exported fields问题
type SerializableWallet struct {
	D         *big.Int
	X, Y      *big.Int
	PublicKey []byte
}

func NewWallets(nodeID string) (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadFromFile(nodeID)

	return &wallets, err
}

func (ws *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := fmt.Sprintf("%s", wallet.GetAddress())

	ws.Wallets[address] = wallet

	return address
}

func (ws *Wallets) GetAddresses() []string {
	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}
func (ws Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

func (ws *Wallets) LoadFromFile(nodeID string) error {

	walletFile := fmt.Sprintf(walletFile, nodeID)
	//fmt.Println(walletFile)
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}

	var wallets map[string]SerializableWallet
	//gob.Register(elliptic.P256())
	gob.Register(SerializableWallet{})
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}
	ws.Wallets = make(map[string]*Wallet)
	//ws.Wallets = wallets.Wallets
	for k, v := range wallets {
		ws.Wallets[k] = &Wallet{
			PrivateKey: ecdsa.PrivateKey{
				PublicKey: ecdsa.PublicKey{
					Curve: elliptic.P256(),
					X:     v.X,
					Y:     v.Y,
				},
				D: v.D,
			},
			PublicKey: v.PublicKey,
		}
	}

	return nil
}

func (ws Wallets) SaveToFile(nodeID string) {
	var content bytes.Buffer
	walletFile := fmt.Sprintf(walletFile, nodeID)
	//fmt.Println(walletFile)

	gob.Register(SerializableWallet{})

	wallets := make(map[string]SerializableWallet)
	for k, v := range ws.Wallets {
		wallets[k] = SerializableWallet{
			D:         v.PrivateKey.D,
			X:         v.PrivateKey.PublicKey.X,
			Y:         v.PrivateKey.PublicKey.Y,
			PublicKey: v.PublicKey,
		}
	}

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(wallets)
	if err != nil {
		log.Panic(err)
	}

	err = os.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}

}
