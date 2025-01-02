// Package contracts contains Solidity contracts and generated Go bindings
package contracts

//go:generate sh -c "solc --bin --abi -o ./bin/ NodeRegistry.sol --overwrite --optimize"
//go:generate sh -c "abigen --bin=./bin/NodeRegistry.bin --abi=./bin/NodeRegistry.abi --pkg=contracts --type NodeRegistry --out=node_registry.go"
