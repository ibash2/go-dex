```shell
npx hardhat help
npx hardhat test
REPORT_GAS=true npx hardhat test
npx hardhat node
npx hardhat run scripts/deploy.js --network gethDev
```

```shell
geth --dev \
     --http \
     --http.addr "0.0.0.0" \
     --http.port 8547 \
     --http.api "eth,net,web3,personal,miner" \
     --http.corsdomain="package://6fd22d6fe5549ad4c4d8fd3ca0b7816b.mod" \
     --ws \
     --ws.addr "0.0.0.0" \
     --ws.port 8546 \
     --ws.api "eth,net,web3,personal,miner" \
     --verbosity 3
```
