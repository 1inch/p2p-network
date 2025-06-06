# Overview
Resolver application serves as a lowest-level endpoint in the p2p-network architecture.

It processes requests received from the relayer and forwards them to the API(s) that it wraps.

Currently 2 APIs are supported: a mock (default) one, and an external Infura API.
```mermaid
---
title: Resolver architecture
---
graph TD
  relayer(Relayer gRPC client node)
  subgraph resolver[Resolver]
    grpc_server[Resolver gRPC server node]

    grpc_server <--> api_wrapper
    subgraph api_wrapper[API handlers]
      default_api[Default]
      infura_api[Infura]
    end
    default_api_impl[Internal mock API impl]
  end
  relayer <==>|gRPC protocol| grpc_server
  infura_api_impl(External Infura API)
  default_api <--> default_api_impl
  infura_api <==>|JSON-RPC| infura_api_impl
```

# Testing notes

## Preparation
First, run the resolver node.

In one terminal session, execute:
```
make build_resolver
bin/resolver run
```
By default resolver listens on port 8001, this can be overridden via `--port` parameter like so:
```
bin/resolver run --port=8888
```
Additionally, one can provide a YAML file with the resolver configuration, and pass its path as the `--config_file` parameter:
```
bin/resolver run --config_file resolver_config.yaml
```

## grpcurl
Now one can test gRPC server responses via `grpcurl`:

1. Failed request (empty JSON)
`grpcurl -plaintext localhost:8001 resolver.Execute/Execute` should return:
```
{
  "error": {
    "message": "empty request id"
  }
}
```
2. Successful request (GetWalletBalance payload):
```
PAYLOAD=$(jq '. | @base64' <<< '{"id": "new", "method": "GetWalletBalance", "params": ["0x1234", "latest"]}')

grpcurl -plaintext -d "{\"id\": \"1\", \"payload\": $PAYLOAD}" localhost:8001 resolver.Execute/Execute | jq '.payload | @base64d | fromjson'

```
Output:
```
{
  "id": "new",
  "result": 555
}
```
**Note**: we first base64-encode the payload to be sent to gRPC. In the above example `jq` is used for base64 encoding.

## Postman
Postman can also be used for testing. 

  - Open Postman
  - In `File->New...` dialog, select `gRPC`.
  - In the Request window, pick your resolver's IP address (by default, this will be `localhost:8001`)
  - And select the service name to test (`Execute/Execute`). If it's not visible, click on `Use Server Reflection`
  - Now, click `Use Example Message` in order to fill in the JSON request template
  - Note that payload has to be base64-encoded, as in the grpcurl example above. The response payload also needs to be decoded.


## HealthCheck
Resolver have standardized health check endpoints:
- Watch - allows you to subscribe to the change service status 
- Check - allows you ask current service status

Both api have same empty request body. You can also try this endpoints in postman using server reflection.


# Metrics
Resolver can return metrics for telemetry, you can go to http endpoint **/metrics** and get metrics.
This option is customizable. You can turn on/off by config.yaml.
For enable the endpoint add in config:
```
metric:
  enabled: true
  port: 8081
```
If you dont want enable the endpoint, dont added this in config or modify the configuration in this way:
```
metric:
  enabled: false
  port: 8081
```
## Parameters
If some parameter doesn`t display on endpoint, you should make a request to resolver and this parameter will be calculated.

- ***grpc_server_connections_open*** [gauge] Number of gRPC server connections open.
- ***grpc_server_connections_total*** [counter] Total number of gRPC server connections opened.
- ***grpc_server_requests_pending{service,method}*** [gauge] Number of gRPC server requests pending.
- ***grpc_server_requests_total{service,method,code}*** [counter] Total number of gRPC server requests completed.
- ***grpc_server_latency_seconds{service,method,code}*** [histogram] Latency of gRPC server requests.
- ***grpc_server_recv_bytes{service,method,frame}*** [histogram] Bytes received in gRPC server requests.
- ***grpc_server_sent_bytes{service,method,frame}***
 [histogram] Bytes sent in gRPC server responses.

## Error Handling

The resolver implements a robust error handling system that covers various failure scenarios:

### Error Types

1. **Configuration Errors**
   - Missing required parameters (e.g., `--config_file`)
   - Invalid configuration values
   - File system errors when loading configuration

2. **gRPC Service Errors**
   ```go
   var (
       errNoHandlerApiInConfig = errors.New("no handler api in config")
       errEmptyRequest = errors.New("empty request")
       errEmptyRequestId = errors.New("empty request id")
       errEmptyPayload = errors.New("empty payload")
       errEmptyPublicKey = errors.New("empty public key")
   )
   ```

3. **Request Processing Errors**
   - Invalid message format
   - Failed response serialization
   - Internal execution failures

### Error Response Format

```protobuf
message Error {
  ErrorCode code = 1;    // Error code
  string message = 2;    // Error description
}

message ResolverResponse {
  string id = 1;
  bool encrypted = 2;
  oneof result {
    bytes payload = 3;  // Successful response
    Error error = 4;    // Error response
  }
}
```

### Error Codes

```go
enum ErrorCode {
  ERR_INVALID_MESSAGE_FORMAT = 0;        // Error in message format
  ERR_GRPC_EXECUTION_FAILED = 1;         // gRPC execution failure
  ERR_RESPONSE_SERIALIZATION_FAILED = 2; // Failed to serialize the response
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
   - Logs errors with appropriate context

3. **Response**
   - Returns structured error responses
   - Includes error codes and messages
   - Maintains request ID correlation

### Example Error Response

```json
{
  "error": {
    "code": 0,
    "message": "Invalid message format: empty request id"
  }
}
```

### Logging

- Errors are logged with appropriate severity levels
- Includes request context for debugging
- Maintains audit trail for troubleshooting
