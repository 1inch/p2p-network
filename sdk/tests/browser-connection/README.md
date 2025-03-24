# Browser connection test

## Overview
This unit tests check creating connection between browser (dapp) and relayer. This tests also added in github CI test job.


## Usages
**You should be located in root of project directory**

1. Run mock discovery service:
```
    docker run -p 8545:8545 --platform linux/amd64 ghcr.io/foundry-rs/foundry:latest "anvil --host 0.0.0.0"
``` 
2. Deploy discovery contract:
```
    make deploy_contract
```
3. Register relayer and resolver:
```
    make register_nodes
```
4. move to dir sdk
```
cd sdk
```
5. install dependencies
```
npm i
```
6. build page for test
```
npm run build
```
7. run tests
```
npm run test
```
