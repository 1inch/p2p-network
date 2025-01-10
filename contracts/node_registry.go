// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// NodeRegistryMetaData contains all meta data concerning the NodeRegistry contract.
var NodeRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"getRelayer\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"ip\",\"type\":\"string\"},{\"internalType\":\"bytes[]\",\"name\":\"publicKeys\",\"type\":\"bytes[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"}],\"name\":\"getResolver\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"ip\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"ip\",\"type\":\"string\"}],\"name\":\"registerRelayer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"ip\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"}],\"name\":\"registerResolver\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b506109268061001c5f395ff3fe608060405234801561000f575f5ffd5b506004361061004a575f3560e01c80636f5a4e381461004e578063bdc5037314610063578063e5ee399814610082578063eea330f914610095575b5f5ffd5b61006161005c36600461057c565b6100b5565b005b61006b610118565b6040516100799291906105e9565b60405180910390f35b61006161009036600461065e565b610278565b6100a86100a336600461057c565b61043e565b60405161007991906106ca565b806101075760405162461bcd60e51b815260206004820152601a60248201527f52656c617965722049502063616e6e6f7420626520656d70747900000000000060448201526064015b60405180910390fd5b5f61011382848361076c565b505050565b6060805f8054610127906106f0565b80601f0160208091040260200160405190810160405280929190818152602001828054610153906106f0565b801561019e5780601f106101755761010080835404028352916020019161019e565b820191905f5260205f20905b81548152906001019060200180831161018157829003601f168201915b505050505091506002805480602002602001604051908101604052809291908181526020015f905b8282101561026e578382905f5260205f200180546101e3906106f0565b80601f016020809104026020016040519081016040528092919081815260200182805461020f906106f0565b801561025a5780601f106102315761010080835404028352916020019161025a565b820191905f5260205f20905b81548152906001019060200180831161023d57829003601f168201915b5050505050815260200190600101906101c6565b5050505090509091565b826102c55760405162461bcd60e51b815260206004820152601b60248201527f5265736f6c7665722049502063616e6e6f7420626520656d707479000000000060448201526064016100fe565b806103125760405162461bcd60e51b815260206004820152601a60248201527f5075626c6963206b65792063616e6e6f7420626520656d70747900000000000060448201526064016100fe565b60018282604051610324929190610826565b908152604051908190036020019020805461033e906106f0565b15905061038d5760405162461bcd60e51b815260206004820152601b60248201527f5265736f6c76657220616c72656164792072656769737465726564000000000060448201526064016100fe565b604080516020601f860181900481028201830183528101858152909182919087908790819085018382808284375f9201919091525050509152506040516001906103da9085908590610826565b908152604051908190036020019020815181906103f79082610835565b5050600280546001810182555f919091527f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace01905061043782848361076c565b5050505050565b60605f60018484604051610453929190610826565b908152604051908190036020019020805461046d906106f0565b80601f0160208091040260200160405190810160405280929190818152602001828054610499906106f0565b80156104e45780601f106104bb576101008083540402835291602001916104e4565b820191905f5260205f20905b8154815290600101906020018083116104c757829003601f168201915b505050505090505f8151116105305760405162461bcd60e51b815260206004820152601260248201527114995cdbdb1d995c881b9bdd08199bdd5b9960721b60448201526064016100fe565b9392505050565b5f5f83601f840112610547575f5ffd5b50813567ffffffffffffffff81111561055e575f5ffd5b602083019150836020828501011115610575575f5ffd5b9250929050565b5f5f6020838503121561058d575f5ffd5b823567ffffffffffffffff8111156105a3575f5ffd5b6105af85828601610537565b90969095509350505050565b5f81518084528060208401602086015e5f602082860101526020601f19601f83011685010191505092915050565b604081525f6105fb60408301856105bb565b828103602084015280845180835260208301915060208160051b840101602087015f5b8381101561065057601f1986840301855261063a8383516105bb565b602095860195909350919091019060010161061e565b509098975050505050505050565b5f5f5f5f60408587031215610671575f5ffd5b843567ffffffffffffffff811115610687575f5ffd5b61069387828801610537565b909550935050602085013567ffffffffffffffff8111156106b2575f5ffd5b6106be87828801610537565b95989497509550505050565b602081525f61053060208301846105bb565b634e487b7160e01b5f52604160045260245ffd5b600181811c9082168061070457607f821691505b60208210810361072257634e487b7160e01b5f52602260045260245ffd5b50919050565b601f82111561011357805f5260205f20601f840160051c8101602085101561074d5750805b601f840160051c820191505b81811015610437575f8155600101610759565b67ffffffffffffffff831115610784576107846106dc565b6107988361079283546106f0565b83610728565b5f601f8411600181146107c9575f85156107b25750838201355b5f19600387901b1c1916600186901b178355610437565b5f83815260208120601f198716915b828110156107f857868501358255602094850194600190920191016107d8565b5086821015610814575f1960f88860031b161c19848701351681555b505060018560011b0183555050505050565b818382375f9101908152919050565b815167ffffffffffffffff81111561084f5761084f6106dc565b6108638161085d84546106f0565b84610728565b6020601f821160018114610895575f831561087e5750848201515b5f19600385901b1c1916600184901b178455610437565b5f84815260208120601f198516915b828110156108c457878501518255602094850194600190920191016108a4565b50848210156108e157868401515f19600387901b60f8161c191681555b50505050600190811b0190555056fea2646970667358221220418e17147bcb61450db72ac911b86f989680179a5aa7e732723f83b0b66358db64736f6c634300081c0033",
}

// NodeRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use NodeRegistryMetaData.ABI instead.
var NodeRegistryABI = NodeRegistryMetaData.ABI

// NodeRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use NodeRegistryMetaData.Bin instead.
var NodeRegistryBin = NodeRegistryMetaData.Bin

// DeployNodeRegistry deploys a new Ethereum contract, binding an instance of NodeRegistry to it.
func DeployNodeRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *NodeRegistry, error) {
	parsed, err := NodeRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(NodeRegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &NodeRegistry{NodeRegistryCaller: NodeRegistryCaller{contract: contract}, NodeRegistryTransactor: NodeRegistryTransactor{contract: contract}, NodeRegistryFilterer: NodeRegistryFilterer{contract: contract}}, nil
}

// NodeRegistry is an auto generated Go binding around an Ethereum contract.
type NodeRegistry struct {
	NodeRegistryCaller     // Read-only binding to the contract
	NodeRegistryTransactor // Write-only binding to the contract
	NodeRegistryFilterer   // Log filterer for contract events
}

// NodeRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type NodeRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodeRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NodeRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodeRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NodeRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodeRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NodeRegistrySession struct {
	Contract     *NodeRegistry     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NodeRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NodeRegistryCallerSession struct {
	Contract *NodeRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// NodeRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NodeRegistryTransactorSession struct {
	Contract     *NodeRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// NodeRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type NodeRegistryRaw struct {
	Contract *NodeRegistry // Generic contract binding to access the raw methods on
}

// NodeRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NodeRegistryCallerRaw struct {
	Contract *NodeRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// NodeRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NodeRegistryTransactorRaw struct {
	Contract *NodeRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNodeRegistry creates a new instance of NodeRegistry, bound to a specific deployed contract.
func NewNodeRegistry(address common.Address, backend bind.ContractBackend) (*NodeRegistry, error) {
	contract, err := bindNodeRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NodeRegistry{NodeRegistryCaller: NodeRegistryCaller{contract: contract}, NodeRegistryTransactor: NodeRegistryTransactor{contract: contract}, NodeRegistryFilterer: NodeRegistryFilterer{contract: contract}}, nil
}

// NewNodeRegistryCaller creates a new read-only instance of NodeRegistry, bound to a specific deployed contract.
func NewNodeRegistryCaller(address common.Address, caller bind.ContractCaller) (*NodeRegistryCaller, error) {
	contract, err := bindNodeRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryCaller{contract: contract}, nil
}

// NewNodeRegistryTransactor creates a new write-only instance of NodeRegistry, bound to a specific deployed contract.
func NewNodeRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*NodeRegistryTransactor, error) {
	contract, err := bindNodeRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryTransactor{contract: contract}, nil
}

// NewNodeRegistryFilterer creates a new log filterer instance of NodeRegistry, bound to a specific deployed contract.
func NewNodeRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*NodeRegistryFilterer, error) {
	contract, err := bindNodeRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryFilterer{contract: contract}, nil
}

// bindNodeRegistry binds a generic wrapper to an already deployed contract.
func bindNodeRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := NodeRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NodeRegistry *NodeRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NodeRegistry.Contract.NodeRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NodeRegistry *NodeRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeRegistry.Contract.NodeRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NodeRegistry *NodeRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NodeRegistry.Contract.NodeRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NodeRegistry *NodeRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NodeRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NodeRegistry *NodeRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NodeRegistry *NodeRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NodeRegistry.Contract.contract.Transact(opts, method, params...)
}

// GetRelayer is a free data retrieval call binding the contract method 0xbdc50373.
//
// Solidity: function getRelayer() view returns(string ip, bytes[] publicKeys)
func (_NodeRegistry *NodeRegistryCaller) GetRelayer(opts *bind.CallOpts) (struct {
	Ip         string
	PublicKeys [][]byte
}, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "getRelayer")

	outstruct := new(struct {
		Ip         string
		PublicKeys [][]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Ip = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.PublicKeys = *abi.ConvertType(out[1], new([][]byte)).(*[][]byte)

	return *outstruct, err

}

// GetRelayer is a free data retrieval call binding the contract method 0xbdc50373.
//
// Solidity: function getRelayer() view returns(string ip, bytes[] publicKeys)
func (_NodeRegistry *NodeRegistrySession) GetRelayer() (struct {
	Ip         string
	PublicKeys [][]byte
}, error) {
	return _NodeRegistry.Contract.GetRelayer(&_NodeRegistry.CallOpts)
}

// GetRelayer is a free data retrieval call binding the contract method 0xbdc50373.
//
// Solidity: function getRelayer() view returns(string ip, bytes[] publicKeys)
func (_NodeRegistry *NodeRegistryCallerSession) GetRelayer() (struct {
	Ip         string
	PublicKeys [][]byte
}, error) {
	return _NodeRegistry.Contract.GetRelayer(&_NodeRegistry.CallOpts)
}

// GetResolver is a free data retrieval call binding the contract method 0xeea330f9.
//
// Solidity: function getResolver(bytes publicKey) view returns(string ip)
func (_NodeRegistry *NodeRegistryCaller) GetResolver(opts *bind.CallOpts, publicKey []byte) (string, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "getResolver", publicKey)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetResolver is a free data retrieval call binding the contract method 0xeea330f9.
//
// Solidity: function getResolver(bytes publicKey) view returns(string ip)
func (_NodeRegistry *NodeRegistrySession) GetResolver(publicKey []byte) (string, error) {
	return _NodeRegistry.Contract.GetResolver(&_NodeRegistry.CallOpts, publicKey)
}

// GetResolver is a free data retrieval call binding the contract method 0xeea330f9.
//
// Solidity: function getResolver(bytes publicKey) view returns(string ip)
func (_NodeRegistry *NodeRegistryCallerSession) GetResolver(publicKey []byte) (string, error) {
	return _NodeRegistry.Contract.GetResolver(&_NodeRegistry.CallOpts, publicKey)
}

// RegisterRelayer is a paid mutator transaction binding the contract method 0x6f5a4e38.
//
// Solidity: function registerRelayer(string ip) returns()
func (_NodeRegistry *NodeRegistryTransactor) RegisterRelayer(opts *bind.TransactOpts, ip string) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "registerRelayer", ip)
}

// RegisterRelayer is a paid mutator transaction binding the contract method 0x6f5a4e38.
//
// Solidity: function registerRelayer(string ip) returns()
func (_NodeRegistry *NodeRegistrySession) RegisterRelayer(ip string) (*types.Transaction, error) {
	return _NodeRegistry.Contract.RegisterRelayer(&_NodeRegistry.TransactOpts, ip)
}

// RegisterRelayer is a paid mutator transaction binding the contract method 0x6f5a4e38.
//
// Solidity: function registerRelayer(string ip) returns()
func (_NodeRegistry *NodeRegistryTransactorSession) RegisterRelayer(ip string) (*types.Transaction, error) {
	return _NodeRegistry.Contract.RegisterRelayer(&_NodeRegistry.TransactOpts, ip)
}

// RegisterResolver is a paid mutator transaction binding the contract method 0xe5ee3998.
//
// Solidity: function registerResolver(string ip, bytes publicKey) returns()
func (_NodeRegistry *NodeRegistryTransactor) RegisterResolver(opts *bind.TransactOpts, ip string, publicKey []byte) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "registerResolver", ip, publicKey)
}

// RegisterResolver is a paid mutator transaction binding the contract method 0xe5ee3998.
//
// Solidity: function registerResolver(string ip, bytes publicKey) returns()
func (_NodeRegistry *NodeRegistrySession) RegisterResolver(ip string, publicKey []byte) (*types.Transaction, error) {
	return _NodeRegistry.Contract.RegisterResolver(&_NodeRegistry.TransactOpts, ip, publicKey)
}

// RegisterResolver is a paid mutator transaction binding the contract method 0xe5ee3998.
//
// Solidity: function registerResolver(string ip, bytes publicKey) returns()
func (_NodeRegistry *NodeRegistryTransactorSession) RegisterResolver(ip string, publicKey []byte) (*types.Transaction, error) {
	return _NodeRegistry.Contract.RegisterResolver(&_NodeRegistry.TransactOpts, ip, publicKey)
}
