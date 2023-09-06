package blockchain

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"recy_network/contrato"
	"recy_network/util"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
)

type Issue struct {
	Client *ethclient.Client
	Config util.Config
}

func NewIssue(config util.Config) *Issue {

	client := CreateClient(config.ProviderCelo)
	if client == nil {
		err := errors.New("Issue:Client NIL")
		log.Println(err)
		os.Exit(1)
	}
	return &Issue{Client: client, Config: config}
}

func (issue *Issue) IssueToken(wallet string, total decimal.Decimal) (string, error) {

	valid := util.IsValidAddress(wallet)
	if !valid {
		return "", fmt.Errorf("wallet INVALID:", wallet)
	}
	address := common.HexToAddress(issue.Config.CRecyAddress)
	instance, err := contrato.NewCrecy(address, issue.Client)
	if err != nil {
		log.Fatal(err)
	}

	secrectKey := issue.Config.SecretCelo
	auth := GetAuth(secrectKey, issue.Client)
	toWallet := common.HexToAddress(wallet)
	totalWei := util.ToWei(total, 18)

	tx, err := instance.Mint(auth, toWallet, totalWei)
	if err != nil {
		log.Fatal(err)
	}

	txSender, _, err := issue.Client.TransactionByHash(context.Background(), tx.Hash())
	receipt, errTx := bind.WaitMined(context.Background(), issue.Client, txSender)
	if errTx != nil {
		log.Println(errTx)
		return "", errTx
	}
	if receipt.Status == types.ReceiptStatusFailed {
		return "", fmt.Errorf(receipt.TxHash.String())
	}
	return tx.Hash().Hex(), nil
}

func (issue *Issue) Close() error {

	if issue.Client == nil {
		return errors.New("Client NIL")
	}

	issue.Client.Close()

	return nil
}
