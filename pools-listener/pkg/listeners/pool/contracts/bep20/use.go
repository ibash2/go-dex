package bep20

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetSymbol(client *ethclient.Client, contractAddress common.Address) string {
	instance, err := NewContracts(contractAddress, client)
	if err != nil {
		return ""
	}
	symbol, err := instance.Symbol(&bind.CallOpts{})
	if err != nil {
		return ""
	}
	return symbol
}

func GetName(client *ethclient.Client, contractAddress common.Address) (string, error) {
	instance, err := NewContracts(contractAddress, client)
	if err != nil {
		return "", err
	}
	name, err := instance.Name(&bind.CallOpts{})
	if err != nil {
		return "", err
	}
	return name, nil
}
