# API description

## Usage
API exports a `Client` class that handles interaction with API backend.

Two methods need to be invoked:

- `init(ClientParams)`: initializes the client by connecting to a discovery service
This discovery service is then used to fetch available relayers and resolvers that will be used for request processing, inside `NetworkParams` object. 

Returns a Promise that is fulfilled once WebRTC connection to relayer is established.
- `execute(JsonRequest): Promise<JsonResponse>`: this actually executes a request and returns a Promise that will contain a `JsonResponse` upon succesful completion, 

Type descriptions below:
```typescript
type NetworkParams = {
  relayerIp: string,
  resolverPubKey: string,
};


type ClientParams = {
  providerUrl: string,
  contractAddr: string,
};

export type JsonRequest = {
  Id: string;
  Method: string;
  Params: string[];
};

export type JsonResponse = {
  id: string;
  result: any;
};
```

Sample code is in `test.ts` file.

## Flow
`Client.init()` performs these actions: 

  - fetches relayer/resolver data from smart contract using `ClientParams`
  - establishes a WebRTC DataChannel connection to the relayer

`Client.execute()` in turn does the following:

- encrypts a request payload with resolve's public key from `NetworkParams` fetched during client initializtion
- generates its own asymmetric Secp256k1 keypair. It is used for response decryption and gets passed alongside request data. It is unique for each request.
- submits protobuf wrapping the request through WebRTC DataChannel
- waits for response on DataChannel's `onmessage` event handler
- upon receiving the response, decrypts it with the corresponding private key

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
