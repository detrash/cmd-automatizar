package main

import (
	"fmt"
	"log"

	"recy_network/blockchain"
	"recy_network/store"
	"recy_network/util"

	"github.com/oklog/ulid/v2"
)

func main() {

	//CertificateMintNFT()

	IssueToken()
}

func CertificateMintNFT() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	var certStore = store.NewDatabaseStore(config)
	strDate := "2022-04-01" //20 ABRIL 2023
	infos, error := certStore.GetInfoForRegistry(strDate)

	if error != nil {
		log.Fatalln(error.Error())
	}

	var registry = blockchain.NewRegistry(config)

	for _, inf := range infos {

		fmt.Println("----------------")
		txHash, err := registry.MintNFTCertificate(inf.FormID, inf.URLCertificate)

		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(txHash)

		nftID, err := registry.GetReturnTxMint(txHash)
		if err != nil {
			log.Fatal(err)
		}
		uuid := ulid.Make()
		inf.ID = uuid.String()
		inf.IDCertificate = nftID
		inf.TxHash = txHash
		inf.FromWallet = config.WalletNFT
		inf.ToWallet = config.WalletNFT

		err = certStore.SaveCertificate(inf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("----------------")
	}

	registry.Close()
	certStore.Close()
}

func IssueToken() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	var dbStore = store.NewDatabaseStore(config)
	issueReports, error := dbStore.GetInfoIssue()

	if error != nil {
		log.Fatalln(error.Error())
	}
	var issue = blockchain.NewIssue(config)

	for _, issueReport := range issueReports {
		uuid := ulid.Make()
		issueReport.ID = uuid.String()
		//WALLET USER
		issueReport.Allocations[0].Wallet = issueReport.Wallet

		var arrAllocation = issueReport.Allocations
		for i, allocation := range arrAllocation {
			txHash, err := issue.IssueToken(allocation.Wallet, allocation.Total)
			if err != nil {
				log.Fatalln(error.Error())
			}
			arrAllocation[i].TxHash = txHash
		}
		dbStore.SaveIssueReport(issueReport)
	}

	issue.Close()
	dbStore.Close()

}
