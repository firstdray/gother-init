package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

func main() {

	entropy, _ := bip39.NewEntropy(128)
	mnemonic, _ := bip39.NewMnemonic(entropy)

	seed := bip39.NewSeed(mnemonic, "")

	masterKey, _ := bip32.NewMasterKey(seed)

	path := []uint32{44 + 0x80000000, 60 + 0x80000000, 0 + 0x80000000, 0, 0}
	key, _ := masterKey.NewChildKey(path[0])
	key, _ = key.NewChildKey(path[1])
	key, _ = key.NewChildKey(path[2])
	key, _ = key.NewChildKey(path[3])
	key, _ = key.NewChildKey(path[4])

	privateKey, err := crypto.ToECDSA(key.Key)
	if err != nil {
		panic(err)
	}

	publicKey := privateKey.Public()

	addr := crypto.PubkeyToAddress(*publicKey.(*ecdsa.PublicKey))
	//addr := common.HexToAddress("0x36449d413789c6Df27D21bD19De3732A1F63CC4C")

	fmt.Println("address:", addr.Hex())

	cl, err := ethclient.Dial("https://endpoints.omniatech.io/v1/bsc/testnet/public")
	if err != nil {
		panic(err)
	}

	balance, err := cl.BalanceAt(context.Background(), addr, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("balance: ", balance)
}
