package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ABI события Swap
const swapEventABI = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"sender","type":"address"},{"indexed":false,"internalType":"uint256","name":"amount0In","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"amount1In","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"amount0Out","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"amount1Out","type":"uint256"},{"indexed":true,"internalType":"address","name":"to","type":"address"}],"name":"Swap","type":"event"}]`

// Адрес контракта
const contractAddress = "0x3e0c0b875002c473047079600c467EB1eD623ea1"

func main() {
	// Подключение к клиенту Ethereum (можно использовать Infura или локальный узел)
	client, err := ethclient.Dial("ws://127.0.0.1:8546")
	if err != nil {
		log.Fatal(err)
	}

	// Преобразование адреса контракта в common.Address
	contractAddress := common.HexToAddress(contractAddress)

	// Разбор ABI контракта для получения информации о событии
	parsedABI, err := abi.JSON(strings.NewReader(swapEventABI))
	if err != nil {
		log.Fatal(err)
	}

	// Получаем хэш события Swap
	swapEventSignature := []byte("Swap(address,uint256,uint256,uint256,uint256,address)")
	swapEventHash := common.BytesToHash(crypto.Keccak256(swapEventSignature))

	// Создаем фильтр для подписки на событие Swap
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		Topics:    [][]common.Hash{{swapEventHash}},
	}

	// Канал для получения логов событий
	logs := make(chan types.Log)

	// Подписка на события
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Listening for Swap events...")

	// Обработка полученных событий
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			// Парсинг данных события
			event := struct {
				Sender     common.Address
				Amount0In  *big.Int
				Amount1In  *big.Int
				Amount0Out *big.Int
				Amount1Out *big.Int
				To         common.Address
			}{}

			err := parsedABI.UnpackIntoInterface(&event, "Swap", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}

			// Вывод информации о событии
			fmt.Printf("Swap Event Detected: Sender=%s, To=%s, Amount0In=%s, Amount1In=%s, Amount0Out=%s, Amount1Out=%s\n",
				event.Sender.Hex(), event.To.Hex(), event.Amount0In.String(), event.Amount1In.String(), event.Amount0Out.String(), event.Amount1Out.String())
		}
	}
}
