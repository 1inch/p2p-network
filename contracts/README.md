# NodeRegistry Deployment Guide

## Table of Contents

- [Run and Deploy](#run-and-deploy)
- [Rebuild Contract](#rebuild-contract)

---

## Run and Deploy

### Prerequisites

- **Anvil**: Ensure Anvil is installed as part of the [Foundry](https://book.getfoundry.sh/) suite.


1. **Start Anvil**

   Launch the local Ethereum node using Anvil:

   ```bash
   make start-anvil
   ```

   To stop:

   ```bash
   make stop-anvil
   ```
      ```

   To stop:

   ```bash
   make stop-anvil
   ```

   Alternatively, you can start Anvil in Docker:

   ```bash
   docker run -p 8545:8545 --platform linux/amd64 ghcr.io/foundry-rs/foundry:latest "anvil --host 0.0.0.0"
   ```

   **Default Anvil Settings:**

   ```markdown
   RPC URL: http://127.0.0.1:8545
   Chain ID: 31337
   Accounts: 10 pre-funded accounts with 1000 ETH each. 
   
   *First account we use to deploy the contract, second is used by the Relayer node, others can be used by the Resolver nodes
   ```


2. **Deploy the Contract**

   Deploy the `NodeRegistry` smart contract to the local Anvil node:

   ```bash
   make deploy_contract
   ```

## Rebuild Contract (optional)

*Use this section if you've made changes to the Solidity contract.*

### Prerequisites

- **Solidity Compiler (`solc`)**: Ensure `solc` is installed. [Solidity Installation Guide](https://docs.soliditylang.org/en/v0.8.17/installing-solidity.html)
- **Go-Ethereum Tools (`abigen`)**: Ensure `abigen` is installed and accessible in your `PATH`.

  **Install `abigen`:**

 ```bash
 make golang_deps
 ```

 *if abigen is not recognized: 
 ```bash
 export PATH=$PATH:$(go env GOPATH)/bin
 ``` 

### Steps

1. **Build the Contract**

   Compile the Solidity contract, generate the ABI and bytecode, generate bindings:

   ```bash
   make generate_bindings
   ```
