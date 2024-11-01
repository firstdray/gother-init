package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

func main() {

	// generating mnemonic
	entropy, _ := bip39.NewEntropy(128)
	mnemonic, _ := bip39.NewMnemonic(entropy)

	// Derive the seed from the mnemonic
	seed := bip39.NewSeed(mnemonic, "")

	// Create a BIP32 master key
	masterKey, _ := bip32.NewMasterKey(seed)

	// Standard derivation path for Ethereum
	path := []uint32{44 + 0x80000000, 60 + 0x80000000, 0 + 0x80000000, 0, 0}

	// Derive the key following the derivation path (sequentially)
	key, _ := masterKey.NewChildKey(path[0])
	key, _ = key.NewChildKey(path[1])
	key, _ = key.NewChildKey(path[2])
	key, _ = key.NewChildKey(path[3])
	key, _ = key.NewChildKey(path[4])

	// Decode the byte array into an ECDSA private key
	privateKey, err := crypto.ToECDSA(key.Key)
	if err != nil {
		panic(err)
	}

	// Get the ECDSA private key
	privateKeyBytes := crypto.FromECDSA(privateKey)
	// Converting a byte array to hex
	privateKeyEth := common.Bytes2Hex(privateKeyBytes)

	// Get the ECDSA public key
	publicKey := privateKey.Public()

	// Calculate the Ethereum address
	addr := crypto.PubkeyToAddress(*publicKey.(*ecdsa.PublicKey))
	//addr := common.HexToAddress("0x36449d413789c6Df27D21bD19De3732A1F63CC4C")

	// Connect to an Ethereum node (replace with the correct URL)
	cl, err := ethclient.Dial("https://endpoints.omniatech.io/v1/bsc/testnet/public")
	if err != nil {
		panic(err)
	}

	defer cl.Close()

	// Get the balance at the address
	balance, err := cl.BalanceAt(context.Background(), addr, nil)
	if err != nil {
		panic(err)
	}

	//tokenAddress := common.HexToAddress("0x1EeB01771dEDa232b46AD9fB637A6D1A7BA1F57B")

	fmt.Println("private key: ", privateKeyEth)
	fmt.Println("address:", addr.Hex())
	fmt.Println("balance: ", balance)
	//fmt.Println("balanceTk: ", balanceOf)
}
