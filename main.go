package main

import (
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"github.com/Zilliqa/gozilliqa-sdk/keytools"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	"github.com/Zilliqa/gozilliqa-sdk/util"
	"strconv"
	"time"
)

func main() {
	// config those parameter probably
	privateKey := ""
	url := "https://dev-api.zilliqa.com/"
	contractAddr := "312253d9bff24384b8dd51d31656b8583de06062"
	contractBech32,_ := bech32.ToBech32Address(contractAddr)
	chainId := 333
	batchNum := 10
	txnsPerBatch := 10

	// no need any change
	msgVersion := 1
	wallet := account.NewWallet()
	wallet.AddByPrivateKey(privateKey)
	client := provider.NewProvider(url)
	pubKey := util.EncodeHex(keytools.GetPublicKeyFromPrivateKey(util.DecodeHex(privateKey), true))

	for i := 0; i < batchNum; i++ {
		fmt.Printf("process batch %d\n", i)
		var transactions []*transaction.Transaction
		for j := 0; j < txnsPerBatch; j++ {
			fmt.Printf("constrct transaction %d\n", j)
			index := fmt.Sprintf("batch %d, num %d", i, j)
			args := []core.ContractValue{
				{
					"key",
					"String",
					"key " + index,
				},
				{
					"value",
					"String",
					"value " + index,
				},
			}
			data := contract2.Data{
				Tag:    "polulate",
				Params: args,
			}

			txn := &transaction.Transaction{
				Version:      strconv.FormatInt(int64(util.Pack(chainId, msgVersion)), 10),
				Amount:       "0",
				GasPrice:     "2000000000",
				GasLimit:     "30000",
				SenderPubKey: pubKey,
				ToAddr:       contractBech32,
				Code:         "",
				Data:         data,
				Priority:     true,
			}
			transactions = append(transactions, txn)

		}
		fmt.Println("sign transactions")
		err := wallet.SignBatch(transactions, *client)
		if err != nil {
			panic(err)
		}
		fmt.Println("send transactions")
		batchSendingResult, err := wallet.SendBatchOneGo(transactions, *client)
		if err != nil {
			panic(err)
		}
		fmt.Println(batchSendingResult)
		fmt.Println("sleep for some time")
		time.Sleep(time.Minute)
	}

}
