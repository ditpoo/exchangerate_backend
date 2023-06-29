package web3

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type Wallet struct {
	Balance string
	Address string
}

func New(address string) *Wallet {
	return &Wallet{ 
		Balance: "",
		Address: address,
	}
}

func (w *Wallet) UpdateBalance() {
	nodeURL := fmt.Sprintf("https://mainnet.infura.io/v3/%s", os.Getenv("INFURA_PROJECT_ID"))
	client, err := rpc.Dial(nodeURL)
	if err != nil {
		log.Fatal(err)
		return
	}
	ethClient := ethclient.NewClient(client)
	parsedAddress := common.HexToAddress(w.Address)
	balance, err := ethClient.BalanceAt(context.Background(), parsedAddress, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	etherBalance := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18))
	w.Balance = etherBalance.String()
	log.Printf("Address: %s\n", parsedAddress.Hex())
	log.Printf("Balance: %s ETH\n", etherBalance.String())
}
