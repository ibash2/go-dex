package pool

import (
	"context"

	"fmt"
	"log"
	"math/big"

	"main/pkg/listeners/pool/contracts/bep20"
	"main/pkg/rmq/event"
	"main/pkg/rmq/publisher"

	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	infuraURL          = "wss://bsc-mainnet.infura.io/ws/v3/583c43516f464b2aae02993372269742"
	contractAddressHex = "0x8341b19a2A602eAE0f22633b6da12E1B016E6451"
	eventSignature     = "Sync(address,uint256,uint256,uint256,uint256,uint256,uint256)"
)

func Run() {
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client: %v", err)
	}

	contractAddress := common.HexToAddress(contractAddressHex)

	// Calculate the event signature hash
	eventSignatureHash := getEventSignatureHash(eventSignature)

	// Create a filter query for logs from the contract
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		Topics:    [][]common.Hash{{eventSignatureHash}},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatalf("Failed to subscribe to logs: %v", err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatalf("Error received from subscription: %v", err)
		case vLog := <-logs:
			// The first topic is the indexed address (the first parameter of Sync)
			topic := vLog.Topics

			// Unpack the data (non-indexed arguments)
			arg1 := common.BytesToAddress(vLog.Data[12:32])
			arg2 := new(big.Int).SetBytes(vLog.Data[32:64])
			arg3 := new(big.Int).SetBytes(vLog.Data[64:96])
			arg4 := new(big.Int).SetBytes(vLog.Data[96:128])
			arg5 := new(big.Int).SetBytes(vLog.Data[128:160])
			arg6 := new(big.Int).SetBytes(vLog.Data[160:192])
			arg7 := new(big.Int).SetBytes(vLog.Data[192:224])

			symbol := bep20.GetSymbol(client, arg1)
			// Output the data
			fmt.Printf("Tx: %s\n", vLog.TxHash.Hex())
			fmt.Printf("Topic: %s\n", topic[0].Hex())
			fmt.Printf("Address: %s [%s]\n", arg1.String(), symbol)
			fmt.Printf("Arg2: %s\n", arg2.String())
			fmt.Printf("Arg3: %s\n", arg3.String())
			fmt.Printf("Arg4: %s\n", arg4.String())
			fmt.Printf("Arg5: %s\n", arg5.String())
			fmt.Printf("Arg6: %s\n", arg6.String())
			fmt.Printf("Arg7: %s\n", arg7.String())

			arg4Float := new(big.Float).SetInt(arg4)
			arg5Float := new(big.Float).SetInt(arg5)
			// difference := new(big.Int).Sub(arg5, arg2)
			// fmt.Printf("Difference: %s\n", difference.String())
			price := new(big.Float).Quo(arg4Float, arg5Float)

			priceFloat, _ := price.Float64()
			fmt.Printf("Price: %s\n\n", price.Text('f', 18))

			publisher.SendEvent(event.NewPriceEvent{Address: arg1.String(), Price: priceFloat})
		}
	}
}

// getEventSignatureHash returns the keccak256 hash of the event signature
func getEventSignatureHash(eventSignature string) common.Hash {
	hash := sha3.NewLegacyKeccak256()
	hash.Write([]byte(eventSignature))
	return common.BytesToHash(hash.Sum(nil))
}
