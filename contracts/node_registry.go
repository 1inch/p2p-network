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
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"ip\",\"type\":\"string\"}],\"name\":\"RelayerRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"ip\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"}],\"name\":\"ResolverRegistered\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"getRelayer\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"ip\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"}],\"name\":\"getResolver\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"ip\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"ip\",\"type\":\"string\"}],\"name\":\"registerRelayer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"ip\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"}],\"name\":\"registerResolver\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b506108a78061001c5f395ff3fe608060405234801561000f575f5ffd5b506004361061004a575f3560e01c80636f5a4e381461004e578063bdc5037314610063578063e5ee399814610081578063eea330f914610094575b5f5ffd5b61006161005c3660046104da565b6100a7565b005b61006b6100fe565b6040516100789190610519565b60405180910390f35b61006161008f36600461054e565b6101dd565b61006b6100a23660046104da565b61037a565b5f6100b3828483610652565b506001805460ff1916811790556040517f8dbe58e1d21d25de93f836ccb277c668a9d6ed4bb2e3c0f56dc2ddd10bd367dc906100f29084908490610734565b60405180910390a15050565b60015460609060ff166101505760405162461bcd60e51b8152602060048201526015602482015274139bc81c995b185e595c881c9959da5cdd195c9959605a1b60448201526064015b60405180910390fd5b5f805461015c906105ce565b80601f0160208091040260200160405190810160405280929190818152602001828054610188906105ce565b80156101d35780601f106101aa576101008083540402835291602001916101d3565b820191905f5260205f20905b8154815290600101906020018083116101b657829003601f168201915b5050505050905090565b600282826040516101ef92919061074f565b9081526040519081900360200190206002015460ff16156102525760405162461bcd60e51b815260206004820152601b60248201527f5265736f6c76657220616c7265616479207265676973746572656400000000006044820152606401610147565b604051806060016040528085858080601f0160208091040260200160405190810160405280939291908181526020018383808284375f92019190915250505090825250604080516020601f8601819004810282018101909252848152918101919085908590819084018382808284375f9201919091525050509082525060016020909101526040516002906102ea908590859061074f565b90815260405190819003602001902081518190610307908261075e565b506020820151600182019061031c908261075e565b50604091820151600291909101805460ff1916911515919091179055517fd05be4c6e5d9c326c84d7bd08078021f3c2ae6decbcf537dfe509225e21192009061036c908690869086908690610819565b60405180910390a150505050565b60606002838360405161038e92919061074f565b9081526040519081900360200190206002015460ff166103e55760405162461bcd60e51b815260206004820152601260248201527114995cdbdb1d995c881b9bdd08199bdd5b9960721b6044820152606401610147565b600283836040516103f792919061074f565b9081526040519081900360200190208054610411906105ce565b80601f016020809104026020016040519081016040528092919081815260200182805461043d906105ce565b80156104885780601f1061045f57610100808354040283529160200191610488565b820191905f5260205f20905b81548152906001019060200180831161046b57829003601f168201915b5050505050905092915050565b5f5f83601f8401126104a5575f5ffd5b50813567ffffffffffffffff8111156104bc575f5ffd5b6020830191508360208285010111156104d3575f5ffd5b9250929050565b5f5f602083850312156104eb575f5ffd5b823567ffffffffffffffff811115610501575f5ffd5b61050d85828601610495565b90969095509350505050565b602081525f82518060208401528060208501604085015e5f604082850101526040601f19601f83011684010191505092915050565b5f5f5f5f60408587031215610561575f5ffd5b843567ffffffffffffffff811115610577575f5ffd5b61058387828801610495565b909550935050602085013567ffffffffffffffff8111156105a2575f5ffd5b6105ae87828801610495565b95989497509550505050565b634e487b7160e01b5f52604160045260245ffd5b600181811c908216806105e257607f821691505b60208210810361060057634e487b7160e01b5f52602260045260245ffd5b50919050565b601f82111561064d57805f5260205f20601f840160051c8101602085101561062b5750805b601f840160051c820191505b8181101561064a575f8155600101610637565b50505b505050565b67ffffffffffffffff83111561066a5761066a6105ba565b61067e8361067883546105ce565b83610606565b5f601f8411600181146106af575f85156106985750838201355b5f19600387901b1c1916600186901b17835561064a565b5f83815260208120601f198716915b828110156106de57868501358255602094850194600190920191016106be565b50868210156106fa575f1960f88860031b161c19848701351681555b505060018560011b0183555050505050565b81835281816020850137505f828201602090810191909152601f909101601f19169091010190565b602081525f61074760208301848661070c565b949350505050565b818382375f9101908152919050565b815167ffffffffffffffff811115610778576107786105ba565b61078c8161078684546105ce565b84610606565b6020601f8211600181146107be575f83156107a75750848201515b5f19600385901b1c1916600184901b17845561064a565b5f84815260208120601f198516915b828110156107ed57878501518255602094850194600190920191016107cd565b508482101561080a57868401515f19600387901b60f8161c191681555b50505050600190811b01905550565b604081525f61082c60408301868861070c565b828103602084015261083f81858761070c565b97965050505050505056fea2646970667358221220cf036eb2630b837d8afab8595db4299d24f7b49402ef818060eb9e4e3853ad8b64736f6c637829302e382e32382d646576656c6f702e323032342e31302e31302b636f6d6d69742e3738393336313461005a",
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
// Solidity: function getRelayer() view returns(string ip)
func (_NodeRegistry *NodeRegistryCaller) GetRelayer(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "getRelayer")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetRelayer is a free data retrieval call binding the contract method 0xbdc50373.
//
// Solidity: function getRelayer() view returns(string ip)
func (_NodeRegistry *NodeRegistrySession) GetRelayer() (string, error) {
	return _NodeRegistry.Contract.GetRelayer(&_NodeRegistry.CallOpts)
}

// GetRelayer is a free data retrieval call binding the contract method 0xbdc50373.
//
// Solidity: function getRelayer() view returns(string ip)
func (_NodeRegistry *NodeRegistryCallerSession) GetRelayer() (string, error) {
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

// NodeRegistryRelayerRegisteredIterator is returned from FilterRelayerRegistered and is used to iterate over the raw logs and unpacked data for RelayerRegistered events raised by the NodeRegistry contract.
type NodeRegistryRelayerRegisteredIterator struct {
	Event *NodeRegistryRelayerRegistered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NodeRegistryRelayerRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryRelayerRegistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NodeRegistryRelayerRegistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NodeRegistryRelayerRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryRelayerRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryRelayerRegistered represents a RelayerRegistered event raised by the NodeRegistry contract.
type NodeRegistryRelayerRegistered struct {
	Ip  string
	Raw types.Log // Blockchain specific contextual infos
}

// FilterRelayerRegistered is a free log retrieval operation binding the contract event 0x8dbe58e1d21d25de93f836ccb277c668a9d6ed4bb2e3c0f56dc2ddd10bd367dc.
//
// Solidity: event RelayerRegistered(string ip)
func (_NodeRegistry *NodeRegistryFilterer) FilterRelayerRegistered(opts *bind.FilterOpts) (*NodeRegistryRelayerRegisteredIterator, error) {

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "RelayerRegistered")
	if err != nil {
		return nil, err
	}
	return &NodeRegistryRelayerRegisteredIterator{contract: _NodeRegistry.contract, event: "RelayerRegistered", logs: logs, sub: sub}, nil
}

// WatchRelayerRegistered is a free log subscription operation binding the contract event 0x8dbe58e1d21d25de93f836ccb277c668a9d6ed4bb2e3c0f56dc2ddd10bd367dc.
//
// Solidity: event RelayerRegistered(string ip)
func (_NodeRegistry *NodeRegistryFilterer) WatchRelayerRegistered(opts *bind.WatchOpts, sink chan<- *NodeRegistryRelayerRegistered) (event.Subscription, error) {

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "RelayerRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryRelayerRegistered)
				if err := _NodeRegistry.contract.UnpackLog(event, "RelayerRegistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRelayerRegistered is a log parse operation binding the contract event 0x8dbe58e1d21d25de93f836ccb277c668a9d6ed4bb2e3c0f56dc2ddd10bd367dc.
//
// Solidity: event RelayerRegistered(string ip)
func (_NodeRegistry *NodeRegistryFilterer) ParseRelayerRegistered(log types.Log) (*NodeRegistryRelayerRegistered, error) {
	event := new(NodeRegistryRelayerRegistered)
	if err := _NodeRegistry.contract.UnpackLog(event, "RelayerRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryResolverRegisteredIterator is returned from FilterResolverRegistered and is used to iterate over the raw logs and unpacked data for ResolverRegistered events raised by the NodeRegistry contract.
type NodeRegistryResolverRegisteredIterator struct {
	Event *NodeRegistryResolverRegistered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NodeRegistryResolverRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryResolverRegistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NodeRegistryResolverRegistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NodeRegistryResolverRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryResolverRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryResolverRegistered represents a ResolverRegistered event raised by the NodeRegistry contract.
type NodeRegistryResolverRegistered struct {
	Ip        string
	PublicKey []byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterResolverRegistered is a free log retrieval operation binding the contract event 0xd05be4c6e5d9c326c84d7bd08078021f3c2ae6decbcf537dfe509225e2119200.
//
// Solidity: event ResolverRegistered(string ip, bytes publicKey)
func (_NodeRegistry *NodeRegistryFilterer) FilterResolverRegistered(opts *bind.FilterOpts) (*NodeRegistryResolverRegisteredIterator, error) {

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "ResolverRegistered")
	if err != nil {
		return nil, err
	}
	return &NodeRegistryResolverRegisteredIterator{contract: _NodeRegistry.contract, event: "ResolverRegistered", logs: logs, sub: sub}, nil
}

// WatchResolverRegistered is a free log subscription operation binding the contract event 0xd05be4c6e5d9c326c84d7bd08078021f3c2ae6decbcf537dfe509225e2119200.
//
// Solidity: event ResolverRegistered(string ip, bytes publicKey)
func (_NodeRegistry *NodeRegistryFilterer) WatchResolverRegistered(opts *bind.WatchOpts, sink chan<- *NodeRegistryResolverRegistered) (event.Subscription, error) {

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "ResolverRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryResolverRegistered)
				if err := _NodeRegistry.contract.UnpackLog(event, "ResolverRegistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseResolverRegistered is a log parse operation binding the contract event 0xd05be4c6e5d9c326c84d7bd08078021f3c2ae6decbcf537dfe509225e2119200.
//
// Solidity: event ResolverRegistered(string ip, bytes publicKey)
func (_NodeRegistry *NodeRegistryFilterer) ParseResolverRegistered(log types.Log) (*NodeRegistryResolverRegistered, error) {
	event := new(NodeRegistryResolverRegistered)
	if err := _NodeRegistry.contract.UnpackLog(event, "ResolverRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
