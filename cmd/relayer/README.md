# Relayer Node

## Overview
The Relayer Node enables clients to interact with the decentralized network. It utilizes HTTP for SDP signaling and WebRTC data channels on the front-facing API while communicating with Resolver nodes through gRPC requests.

## Configuration

The Relayer Node is configured using a YAML file. The example configuration file is in the root of the repo.
Below is an example configuration:

```yaml
log_level: DEBUG
http_endpoint: 127.0.0.1:8880
private_key: 59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d
discovery:
  rpc_url: http://127.0.0.1:8545
  with_node_registry: false
  contract_address: 0x5fbdb2315678afecb367f032d93f642f64180aa3
webrtc:
  ice_server: stun:stun1.l.google.com:19302
  retry:
    enabled: true
    count: 5
    interval: 1s
  port:
    enabled: true
    min: 15000
    max: 15500
```



### Configuration Fields
- **`log_level`**: The logging level for the node (`DEBUG`, `INFO`, `WARN`, `ERROR`).
- **`http_endpoint`**: The HTTP endpoint where the node listens for SDP signaling requests.
- **`private_key`**: The private key belong to the relayer.
- **`discovery.rpc_url`**:  The rpc endpoint of discovery service, expect ETH blockchain node.
- **`discovery.contract_address`**: The address where discovery contract is located.
- **`webrtc.ice_server`**: The ICE server used for WebRTC signaling (e.g., STUN or TURN server).
- **`webrtc.retry.enabled`**: The flag for turn on/off retry request if resolver return some error.
- **`webrtc.retry.count`**: The count of attempt repeated requests.
- **`webrtc.retry.interval`**: The interval between repeated requests
- **`webrtc.port.enabled`**: The flag for turn on/off range for peer connections port
- **`webrtc.port.min`**: The minimum from range
- **`webrtc.port.max`**: The maximum from range


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
