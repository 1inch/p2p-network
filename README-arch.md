# Architecture overview

There are several key components to the architecture:

- dApp SDK
- Relayers
- Resolvers
- Discovery service
- Monitoring services

```mermaid
---
title: Architecture overview
---
flowchart TD
    sdk[dAPP SDK] <-->|WebRTC DataChannel| relayer(Relayer)
    sdk -->|getRelayer| discovery(Discovery service)
    relayer -->|registerRelayer| discovery
    relayer -->|getResolver| discovery
    relayer <-->|Execute using gRPC| resolver(Resolver)
    resolver -->|registerResolver| discovery
    metrics[Metrics service]
```

The successful flow is described in the below diagram:
```mermaid
---
title: Basic flow
---
sequenceDiagram
    participant sdk as dApp SDK
    participant discovery as Discovery service
    participant relayer as Relayer
    participant resolver as Resolver

    relayer->>discovery: registerRelayer
    discovery-->>relayer: success
    resolver->>discovery: registerResolver
    discovery-->>resolver: success
    sdk->>discovery: getRelayer
    discovery-->>sdk: relayer IP + resolver public Keys
    sdk->>relayer: ExecuteRequest (WebRTC)
    relayer->>discovery: getResolver
    discovery-->>relayer: resolver's IP
    relayer->>resolver: ExecuteRequest (gRPC)
    resolver-->>relayer: ExecuteResponse (gRPC)
    relayer-->>sdk: ExecuteResponse (WebRTC)
```

## Payloads
Payloads travelling between dApp SDK, Relayers and Resolvers are Protobuf-wrapped JSON-RPC requests, with the following structure:
```mermaid
---
title: Payload structure

---
classDiagram
    class ResolverRequest {
        <<protobuf>>
        id: string
        payload: bytes
        encrypted: bool
        publicKey: bytes
    }
    class JsonRequest {
        <<json>>
        id: string
        method: string
        params: []string
    }
    class JsonResponse {
        <<json>>
        id: string
        result: bytes
    }
    class ResolverResponse {
        <<protobuf>>
        id: string
        payload: bytes
        encrypted: bool
    }
    class IncomingMessage {
        <<protobuf>>
        publicKeys: bytes
        request: ResolverRequest
    }
    class OutgoingMessage {
        <<protobuf>>
        response: ResolverResponse
        error: Error
    }
    JsonRequest --o ResolverRequest
    JsonResponse --o ResolverResponse
    ResolverRequest --o IncomingMessage
    ResolverResponse --o OutgoingMessage
```

## dApp SDK
dApp SDK is a Typescript library that provides the following functionality:

- request execution
- encryption

SDK communicates with Relayers using WebRTC suite of protocols.

Prior to establishing communication with the Relayer, SDK asks Discovery service for a Relayer node candidate to use.

## Relayers

Relayers perform the following:

- register themselves via the Discovery Service
- process incoming WebRTC requests from dApp SDK
- pass request payloads to Resolvers via gRPC protocol

Relayers act as gRPC clients in the p2p-network architecture.

More info is in a dedicated [README](./cmd/relayer/README.md).

## Resolvers
Resolvers implement APIs supported by the p2p-network. 

They act as gRPC servers in the overall p2p-network architecture.

More info is in a dedicated [README](./cmd/resolver/README.md).


## Discovery service
Discovery service is a Ethereum smart contract that provides the following functionality:

- relayer registration (**registerRelayer(ip)**)
- resolver registration (**registerResolver(ip, pubKey)**)
- getting relayer and resolver public keys (**getRelayer()**)
- fetching resolver IPs by public key (**getResolver(pubKey)**)



## Monitoring services

# Deployment
All application components are Dockerized. There is a `docker-compose.yml` that can be used for containerized deployment of the architecture described above (excluding monitoring services).
