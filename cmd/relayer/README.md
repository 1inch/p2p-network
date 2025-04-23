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
- **`webrtc.ice_servers.url`**: The ICE server used for WebRTC signaling (e.g., STUN or TURN url server).
- **`webrtc.ice_servers.username`**: The username for TURN server.
- **`webrtc.ice_servers.password`**: The password for TURN server.
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

## Error codes

```
// Enum to represent standardized error codes.
enum ErrorCode {
  ERR_INVALID_MESSAGE_FORMAT = 0;    // Error in message format.
  ERR_RESOLVER_LOOKUP_FAILED = 1;    // Failed to resolve address for public key.
  ERR_GRPC_EXECUTION_FAILED = 2;     // gRPC execution failure.
  ERR_RESPONSE_SERIALIZATION_FAILED = 3; // Failed to serialize the response.
  ERR_DATA_CHANNEL_SEND_FAILED = 4;  // Failed to send the response via the data channel.
}
```

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

## Error Handling

The relayer implements a comprehensive error handling system that covers various failure scenarios:

### Error Types

1. **Configuration Errors**
   - Missing required parameters (e.g., `--config`)
   - Invalid configuration values
   - File system errors when loading configuration

2. **gRPC Service Errors**
   ```go
   var (
       ErrResolverLookupFailed = errors.New("resolver lookup failed")
       ErrGRPCExecutionFailed = errors.New("gRPC execution failed")
       ErrGRPCConnectionCloseFailed = errors.New("gRPC connection close failed")
   )
   ```

3. **Request Processing Errors**
   - Invalid message format
   - Failed response serialization
   - Internal execution failures
   - Network communication errors

### Error Response Format

```protobuf
message Error {
  ErrorCode code = 1;    // Error code
  string message = 2;    // Error description
}

message OutgoingMessage {
  oneof result {
    ResolverResponse response = 1;  // Successful response
    Error error = 2;               // Error response
  }
  bytes publicKey = 3;
}
```

### Error Codes

```go
enum ErrorCode {
  ERR_INVALID_MESSAGE_FORMAT = 0;        // Error in message format
  ERR_RESOLVER_LOOKUP_FAILED = 1;        // Failed to resolve address for public key
  ERR_GRPC_EXECUTION_FAILED = 2;         // gRPC execution failure
  ERR_RESPONSE_SERIALIZATION_FAILED = 3; // Failed to serialize the response
  ERR_DATA_CHANNEL_SEND_FAILED = 4;      // Failed to send the response via the data channel
}
```

### Error Handling Flow

1. **Request Validation**
   - Validates incoming request format
   - Checks for required fields
   - Verifies message integrity

2. **Processing**
   - Handles gRPC execution errors
   - Manages serialization failures
   - Handles network communication errors
   - Implements timeout handling

3. **Response**
   - Returns structured error responses
   - Includes error codes and messages
   - Maintains request ID correlation

### Example Error Response

```json
{
  "error": {
    "code": 1,
    "message": "Failed to resolve address for public key: resolver lookup failed"
  }
}
```

### Logging

- Errors are logged with appropriate severity levels
- Includes request context for debugging
- Maintains audit trail for troubleshooting
- Implements structured logging for better error analysis

### Retry Mechanism

The relayer implements a retry mechanism for transient failures:
- Configurable retry attempts
- Exponential backoff
- Maximum retry delay
- Retryable error types

```yaml
webrtc:
  retry:
    enabled: true
    count: 5
    interval: 1s
```

## Metrics

The relayer provides comprehensive metrics for monitoring and telemetry. These metrics are exposed via a Prometheus-compatible HTTP endpoint at `/metrics`.

### Available Metrics

The relayer exposes the following metrics:

#### HTTP Metrics
- **`relayer_http_requests_total`** [counter] - Total number of HTTP requests received, labeled by handler, method, and status
- **`relayer_http_request_duration_seconds`** [histogram] - Duration of HTTP requests in seconds, labeled by handler and method

#### WebRTC Connection Metrics
- **`relayer_active_peer_connections`** [gauge] - Current number of active PeerConnections
- **`relayer_sdp_negotiation_total`** [counter] - Total number of SDP negotiations, labeled by status
- **`relayer_sdp_negotiation_duration_seconds`** [histogram] - Duration of SDP negotiations in seconds

#### ICE Candidate Metrics
- **`relayer_ice_candidate_sent_total`** [counter] - Total number of ICE candidates sent, labeled by session_id and status
- **`relayer_ice_candidate_send_duration_seconds`** [histogram] - Duration for sending ICE candidates in seconds, labeled by session_id

#### gRPC Metrics
- **`relayer_grpc_requests_total`** [counter] - Total number of gRPC requests, labeled by method and status
- **`relayer_grpc_request_duration_seconds`** [histogram] - Duration of gRPC requests in seconds, labeled by method

#### Data Channel Metrics
- **`relayer_data_channel_messages_sent_total`** [counter] - Total number of messages sent over data channels, labeled by session_id and status
- **`relayer_data_channel_messages_received_total`** [counter] - Total number of messages received over data channels, labeled by session_id
- **`relayer_data_channel_latency_seconds`** [histogram] - Time taken to process data channel messages in seconds, labeled by session_id

### Accessing Metrics

The metrics endpoint is available at the same HTTP server as the other relayer endpoints:

```
http://<relayer-address>:<port>/metrics
```

### Prometheus Integration

The relayer metrics can be integrated with Prometheus by adding the following configuration to your Prometheus setup:

```yaml
scrape_configs:
  - job_name: 'relayer'
    static_configs:
      - targets: ['relayer:8080']
```

### Grafana Dashboard

A Grafana dashboard is available for visualizing the relayer metrics. The dashboard includes panels for:

- Active peer connections
- SDP negotiation success/failure rates
- HTTP request rates and latencies
- gRPC request rates and latencies
- Data channel message rates and latencies
