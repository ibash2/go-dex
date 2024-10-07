package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Connect to the Ethereum client
	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// Retrieve the private key and create a new transactor
	privateKey, err := crypto.HexToECDSA("YOUR_PRIVATE_KEY")
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalf("Failed to cast public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("Failed to get nonce: %v", err)
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get chain ID: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to get gas price: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatalf("Failed to create transactor: %v", err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) // in wei
	auth.GasLimit = uint64(3000000)
	auth.GasPrice = gasPrice

	fmt.Printf("Deploying contracts with the account: %s\n", fromAddress.Hex())

	// Deploy Factory contract
	factoryAddress, tx, _, err := UniswapV2Factory.DeployUniswapV2Factory(auth, client, fromAddress)
	if err != nil {
		log.Fatalf("Failed to deploy factory contract: %v", err)
	}
	fmt.Printf("Factory deployed to %s\n", factoryAddress.Hex())

	// Deploy USDT contract
	usdtAddress, tx, _, err := DeployUSDT(auth, client)
	if err != nil {
		log.Fatalf("Failed to deploy USDT contract: %v", err)
	}
	fmt.Printf("USDT deployed to %s\n", usdtAddress.Hex())

	// Deploy USDC contract
	usdcAddress, tx, _, err := DeployUSDC(auth, client)
	if err != nil {
		log.Fatalf("Failed to deploy USDC contract: %v", err)
	}
	fmt.Printf("USDC deployed to %s\n", usdcAddress.Hex())

	// Mint tokens to the owner
	usdtInstance, err := NewUSDT(usdtAddress, client)
	if err != nil {
		log.Fatalf("Failed to create USDT instance: %v", err)
	}

	tx, err = usdtInstance.Mint(auth, fromAddress, big.NewInt(1000))
	if err != nil {
		log.Fatalf("Failed to mint USDT tokens: %v", err)
	}
	fmt.Printf("Minted 1000 USDT tokens to %s\n", fromAddress.Hex())

	usdcInstance, err := NewUSDC(usdcAddress, client)
	if err != nil {
		log.Fatalf("Failed to create USDC instance: %v", err)
	}

	tx, err = usdcInstance.Mint(auth, fromAddress, big.NewInt(1000))
	if err != nil {
		log.Fatalf("Failed to mint USDC tokens: %v", err)
	}
	fmt.Printf("Minted 1000 USDC tokens to %s\n", fromAddress.Hex())

	// Create pair
	factoryInstance, err := UniswapV2Factory.NewUniswapV2Factory(factoryAddress, client)
	if err != nil {
		log.Fatalf("Failed to create factory instance: %v", err)
	}

	tx, err = factoryInstance.CreatePair(auth, usdtAddress, usdcAddress)
	if err != nil {
		log.Fatalf("Failed to create pair: %v", err)
	}

	pairAddress, err := factoryInstance.GetPair(nil, usdtAddress, usdcAddress)
	if err != nil {
		log.Fatalf("Failed to get pair address: %v", err)
	}
	fmt.Printf("Pair deployed to %s\n", pairAddress.Hex())

	pairInstance, err := IUniswapV2Pair.NewIUniswapV2Pair(pairAddress, client)
	if err != nil {
		log.Fatalf("Failed to create pair instance: %v", err)
	}

	reserves, err := pairInstance.GetReserves(nil)
	if err != nil {
		log.Fatalf("Failed to get reserves: %v", err)
	}
	fmt.Printf("Reserves: %s, %s\n", reserves.Reserve0.String(), reserves.Reserve1.String())

	// Deploy WETH contract
	wethAddress, tx, _, err := DeployWETH(auth, client)
	if err != nil {
		log.Fatalf("Failed to deploy WETH contract: %v", err)
	}
	fmt.Printf("WETH deployed to %s\n", wethAddress.Hex())

	// Deploy Router contract
	routerAddress, tx, _, err := UniswapV2Router02.DeployUniswapV2Router02(auth, client, factoryAddress, wethAddress)
	if err != nil {
		log.Fatalf("Failed to deploy router contract: %v", err)
	}
	fmt.Printf("Router deployed to %s\n", routerAddress.Hex())

	// Approve tokens for router
	maxUint256 := new(big.Int).SetBytes([]byte("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"))

	tx, err = usdtInstance.Approve(auth, routerAddress, maxUint256)
	if err != nil {
		log.Fatalf("Failed to approve USDT tokens: %v", err)
	}

	tx, err = usdcInstance.Approve(auth, routerAddress, maxUint256)
	if err != nil {
		log.Fatalf("Failed to approve USDC tokens: %v", err)
	}

	// Add liquidity
	routerInstance, err := UniswapV2Router02.NewUniswapV2Router02(routerAddress, client)
	if err != nil {
		log.Fatalf("Failed to create router instance: %v", err)
	}

	token0Amount := big.NewInt(100)
	token1Amount := big.NewInt(100)

	deadline := big.NewInt(time.Now().Add(time.Minute * 10).Unix())

	tx, err = routerInstance.AddLiquidity(auth, usdtAddress, usdcAddress, token0Amount, token1Amount, big.NewInt(0), big.NewInt(0), fromAddress, deadline)
	if err != nil {
		log.Fatalf("Failed to add liquidity: %v", err)
	}

	// Check LP token balance for the owner
	lpTokenBalance, err := pairInstance.BalanceOf(nil, fromAddress)
	if err != nil {
		log.Fatalf("Failed to get LP token balance: %v", err)
	}
	fmt.Printf("LP tokens for the owner: %s\n", lpTokenBalance.String())

	reserves, err = pairInstance.GetReserves(nil)
	if err != nil {
		log.Fatalf("Failed to get reserves: %v", err)
	}
	fmt.Printf("Reserves: %s, %s\n", reserves.Reserve0.String(), reserves.Reserve1.String())

	fmt.Printf("USDT_ADDRESS: %s\n", usdtAddress.Hex())
	fmt.Printf("USDC_ADDRESS: %s\n", usdcAddress.Hex())
	fmt.Printf("WETH_ADDRESS: %s\n", wethAddress.Hex())
	fmt.Printf("FACTORY_ADDRESS: %s\n", factoryAddress.Hex())
	fmt.Printf("ROUTER_ADDRESS: %s\n", routerAddress.Hex())
	fmt.Printf("PAIR_ADDRESS: %s\n", pairAddress.Hex())
}
