package blockchain

import (
	"context"
	"errors"
	"log"
	"math/big"
	"os"
	"recy_network/contrato"
	"recy_network/util"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Registry struct {
	Client *ethclient.Client
	Config util.Config
}

func NewRegistry(config util.Config) *Registry {

	client := CreateClient(config.ProviderPol)
	if client == nil {
		err := errors.New("Registry:Client NIL")
		log.Println(err)
		os.Exit(1)
	}
	return &Registry{Client: client, Config: config}
}

func (registry *Registry) MintNFTCertificate(formID, url string) (string, error) {

	contractAddress := registry.Config.NFTAddress
	toAddress := registry.Config.WalletNFT
	address := common.HexToAddress(contractAddress)
	instance, err := contrato.NewContrato(address, registry.Client)
	if err != nil {
		log.Fatal(err)
	}
	//recipient common.Address, tokenURI string
	toMintAddress := common.HexToAddress(toAddress)
	secrectKey := registry.Config.SecretPol
	auth := GetAuth(secrectKey, registry.Client)
	tx, err := instance.MintNFT(auth, toMintAddress, url)
	if err != nil {
		log.Fatal(err)
	}

	return tx.Hash().Hex(), nil

}

func (registry *Registry) Close() error {

	if registry.Client == nil {
		return errors.New("Client NIL")
	}

	registry.Client.Close()

	return nil
}

func (registry *Registry) GetReturnTxMint(strTxHash string) (uint64, error) {

	var nftID uint64 = 0
	if registry.Client == nil {
		return nftID, errors.New("Client NIL")
	}

	txHash := common.HexToHash(strTxHash)

	txSender, _, err := registry.Client.TransactionByHash(context.Background(), txHash)

	receipt, err := bind.WaitMined(context.Background(), registry.Client, txSender)
	if err != nil {
		log.Fatal(err)
	}

	logID := registry.Config.NFTAddress
	for _, vLog := range receipt.Logs {

		if vLog.Address.String() == logID {
			if len(vLog.Topics) > 2 {
				id := new(big.Int)
				id.SetBytes(vLog.Topics[3].Bytes())
				nftID = id.Uint64()
			}
		}
	}

	return nftID, nil
}
