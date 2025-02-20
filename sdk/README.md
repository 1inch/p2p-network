# How to test


## JavaScript side

1. Run the Node express server:

```
cd sdk
npm i
node index.js
```

2. Compile TypeScript files. Below an example using Bun:

```
bun build test.ts --outdir ./build --sourcemap=external --watch
```

## Contract deployment
Run Anvil e.g. via Docker:
```
docker run -p 8545:8545 --platform linux/amd64 ghcr.io/foundry-rs/foundry:latest "anvil --host 0.0.0.0"
```

Then issue the following commands in order to deploy mock registry contract and register default test resolver:
```
make deploy_contract
make register_resolver
```

## Run relayer node

```
make build_relayer_local
bin/relayer run --config relayer.config.example.yaml
```

## Run resolver node
```
make build_resolver_local
bin/resolver run --api=infura --infuraKey=a8401733346d412389d762b5a63b0bcf --privateKey=5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a  --grpc_endpoint=127.0.0.1:8001
```

# Execution
Navigate browser to `http://localhost:3000/index.html` and click `Test Execute`.
