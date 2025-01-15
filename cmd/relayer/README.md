# Relayer Node

## Overview
The Relayer Node enables clients to interact with the decentralized network. It utilizes HTTP for SDP signaling and WebRTC data channels on the front-facing API while communicating with Resolver nodes through gRPC requests.

## Configuration

The Relayer Node is configured using a YAML file. The example configuration file is in the root of the repo.
Below is an example configuration:

```yaml
log_level: DEBUG
http_endpoint: 127.0.0.1:8080
webrtc_ice_server: stun:stun1.l.google.com:19302
grpc_server_address: 127.0.0.1:0
```


```markdown
### Configuration Fields
- **`log_level`**: The logging level for the node (`DEBUG`, `INFO`, `WARN`, `ERROR`).
- **`http_endpoint`**: The HTTP endpoint where the node listens for SDP signaling requests.
- **`webrtc_ice_server`**: The ICE server used for WebRTC signaling (e.g., STUN or TURN server).
- **`grpc_server_address`**: The gRPC server address the node interacts with (Resolver node).
```

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
