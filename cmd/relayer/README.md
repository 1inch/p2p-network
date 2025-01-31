# Relayer Node

## Overview
The Relayer Node enables clients to interact with the decentralized network. It utilizes HTTP for SDP signaling and WebRTC data channels on the front-facing API while communicating with Resolver nodes through gRPC requests.

## Configuration

The Relayer Node is configured using a YAML file. The example configuration file is in the root of the repo.
Below is an example configuration:

```yaml
log_level: DEBUG
http_endpoint: 127.0.0.1:8880
webrtc_ice_server: stun:stun1.l.google.com:19302
grpc_server_address: 127.0.0.1:0
blockchain_rpc_address: http://127.0.0.1:8545
with_node_registry: false
contract_address: 0x5fbdb2315678afecb367f032d93f642f64180aa3
private_key: 59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d
retry_request:
  enabled: true
  count: 5
  interval: 1s
```



### Configuration Fields
- **`log_level`**: The logging level for the node (`DEBUG`, `INFO`, `WARN`, `ERROR`).
- **`http_endpoint`**: The HTTP endpoint where the node listens for SDP signaling requests.
- **`webrtc_ice_server`**: The ICE server used for WebRTC signaling (e.g., STUN or TURN server).
- **`grpc_server_address`**: The gRPC server address the node interacts with (Resolver node).
- **`blockchain_rpc_address`**: The rpc endpoint of ETH blockchain node.
- **`contract_address`**: The address where discovery contract is located.
- **`private_key`**: The private key belong to the relayer.
- **`retry_request.enabled`**: The flag for turn on/off retry request if resolver return some error.
- **`retry_request.count`**: The count of attempt repeated requests.
- **`retry_request.interval`**: The interval between repeated requests


## Command-Line Interface

The Relayer Node is managed via a CLI. Below is the structure of the CLI:

### Commands

- **`run`**: Starts the Relayer Node.
  - **Flags**:
    - `--config`: Path to the YAML configuration file (required).

### Example Usage

To start the Relayer Node with a configuration file:

```bash
./bin/relayer run --config=path/to/config.yaml
```

## Testing the Relayer Node

Currently, there is no SDK available for interacting with the Relayer Node. All testing is performed through the test cases provided in the codebase. These tests validate the functionality of the HTTP API, WebRTC signaling, and data channel communication, as well as the integration with gRPC services.

To run the tests, use the following command:

```bash
make test
```


## HealthCheck
Relayer have http health check endpoint:
- **/health** - allows you ask current service status

The api have empty request body. If service is ok, the api return http status **200 ok**.
