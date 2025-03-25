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
4. start relayer
```
make run_relayer_local
```
5. start resolver
```
make run_resolver_local
```
6. move to directory with tests
```
cd sdk/tests/browser-connection
```
7. install dependencies
```
npm install
```
8. build test script
```
bun build src/call_client_script.ts --outdir ./build --sourcemap=external
```
9. run tests
```
npm test
```
