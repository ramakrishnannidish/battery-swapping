# Install and Run DLT Component
## Download and install Hyperledger Fabric Binaries
This would be required if nvm path is not exported 
```bash
export NVM_DIR="$([ -z "${XDG_CONFIG_HOME-}" ] && printf %s "${HOME}/.nvm" || printf %s "${XDG_CONFIG_HOME}/nvm")"[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh" # This loads nvm
```

```bash
curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && chmod +x install-fabric.sh
./install-fabric.sh docker samples
./install-fabric.sh --fabric-version 2.4.7 binary
```

## Run the network
```bash
cd fabric-samples/test-network
./network.sh up
```

## Create the channel
```bash
./network.sh createChannel
```

## Deploy the Chaincode
Copy the chaincode and application inside the fabric-samples directory for easier path handling (absolute path might required modification)
```bash
cp -R ../../battery-swapping-basic ../
./network.sh deployCC -ccn basic -ccp ../battery-swapping-basic/chaincode-go -ccl go
```


# Run Simulation Application and Dashboard
## Install and run the Simulation Application
```bash
cd ./fabric-samples/battery-swapping-basic/application-gateway-javascript
npm install && npm start
```

## Install and run the Blockchain Explorer
```bash
cd ./fabric-samples/battery-swapping-basic/explorer
```

Need to copy the peer organization crypto-artifacts
```bash
cp -R ../../test-network/organizations ./ 
```

Run the exporer docker containers
```bash
docker-compose up -d
```