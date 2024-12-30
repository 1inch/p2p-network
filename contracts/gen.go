// Package contracts contains Solidity contracts and generated Go bindings
package contracts

//go:generate sh -c "docker run --rm -v $(pwd):/sources ethereum/solc:stable -o /sources/output --abi --bin /sources/NodeRegistry.sol --overwrite --optimize"
//go:generate sh -c "abigen --bin=./output/NodeRegistry.bin --abi=./output/NodeRegistry.abi --pkg=contracts --type NodeRegistry --out=node_registry.go"
