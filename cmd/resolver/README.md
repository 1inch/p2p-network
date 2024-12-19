# Testing notes

## grpcurl
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

Now one can test gRPC server responses via `grpcurl`:

1. Failed request (empty JSON)
`grpcurl -plaintext localhost:8001 proto.Execute/Execute` should return:
```
{
  "status": "RESOLVER_ERROR"
}
```
2. Successful request (GetWalletBalance payload):
```
PAYLOAD=$(jq '. | @base64' <<< '{"id": "new", "method": "GetWalletBalance", "params": ["0x1234", "latest"]}')

grpcurl -plaintext -d "{\"id\": \"1\", \"payload\": $PAYLOAD}" localhost:8001 proto.Execute/Execute | jq '.payload | @base64d | fromjson'

```
Output:
```
{
  "id": "new",
  "result": 0,
  "error": null
}
```
