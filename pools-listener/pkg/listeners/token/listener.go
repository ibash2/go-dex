package token

import (
	"context"
	"fmt"
	"log"

	"math/big"
	"strings"

	"main/pkg/rmq/event"
	"main/pkg/rmq/publisher"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	infuraURL    = "wss://bsc-mainnet.infura.io/ws/v3/583c43516f464b2aae02993372269742"
	contractAddr = "0x0BFbCF9fa4f9C56B0F40a671Ad40E0805A091865"
	eventABI     = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"token0","type":"address"},{"indexed":true,"internalType":"address","name":"token1","type":"address"},{"indexed":true,"internalType":"uint24","name":"fee","type":"uint24"},{"indexed":false,"internalType":"int24","name":"tickSpacing","type":"int24"},{"indexed":false,"internalType":"address","name":"pool","type":"address"}],"name":"PoolCreated","type":"event"}]`
)

func Run() {
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	contractAddress := common.HexToAddress(contractAddr)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatalf("Failed to subscribe to logs: %v", err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(eventABI))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatalf("Error received from subscription: %v", err)
		case vLog := <-logs:
			// Create a struct for the contract_event data (non-indexed parameters)
			contract_event := struct {
				TickSpacing *big.Int
				Pool        common.Address
			}{}

			// Unpack the non-indexed parameters from the log data
			err := contractAbi.UnpackIntoInterface(&contract_event, "PoolCreated", vLog.Data)
			if err != nil {
				log.Fatalf("Failed to unpack log data: %v", err)
			}

			// Decode the indexed parameters from the log topics
			token0 := common.HexToAddress(vLog.Topics[1].Hex())
			token1 := common.HexToAddress(vLog.Topics[2].Hex())
			fee := new(big.Int).SetBytes(vLog.Topics[3].Bytes()) // Indexed fee is stored as a topic

			fmt.Printf("Token0: %s\n", token0.Hex())
			fmt.Printf("Token1: %s\n", token1.Hex())
			fmt.Printf("Fee: %d\n", fee)
			fmt.Printf("TickSpacing: %d\n", contract_event.TickSpacing)
			fmt.Printf("Pool: %s\n\n", contract_event.Pool.Hex())

			publisher.SendEvent(
				event.NewTokenEvent{
					Address:  token0.Hex(),
					Address2: token1.Hex(),
					Chain:    "bsc",
				})
		}
	}
}
