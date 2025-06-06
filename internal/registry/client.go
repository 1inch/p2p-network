// Package registry provides client implementation for node registry.
package registry

import (
	"context"
	"errors"
	"time"

	"github.com/1inch/p2p-network/contracts"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	// ErrContextCancelled error represents context cancellation.
	ErrContextCancelled = errors.New("context cancelled")
	// ErrTransactionFailed error represents transaction failure.
	ErrTransactionFailed = errors.New("transaction failed")
)

// Config represents registry client config.
type Config struct {
	DialURI         string
	PrivateKey      string
	ContractAddress string
}

// Client represents storage client.
type Client struct {
	Registry *contracts.NodeRegistry
	Auth     *bind.TransactOpts
	client   *ethclient.Client
	ticker   *time.Ticker
}

// Dial creates eth client, new smart-contract instance, auth.
func Dial(ctx context.Context, config *Config) (*Client, error) {
	client, err := ethclient.Dial(config.DialURI)
	if err != nil {
		return &Client{}, err
	}

	privateKey, err := crypto.HexToECDSA(config.PrivateKey)
	if err != nil {
		return &Client{}, err
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return &Client{}, err
	}

	registry, err := contracts.NewNodeRegistry(common.HexToAddress(config.ContractAddress), client)
	if err != nil {
		return &Client{}, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return &Client{}, err
	}

	return &Client{
		Registry: registry,
		Auth:     auth,
		client:   client,
		ticker:   time.NewTicker(200 * time.Millisecond),
	}, nil
}

// DeployNodeRegistry deploys node registry smart contract and returns it's client.
func DeployNodeRegistry(ctx context.Context, config *Config) (common.Address, *Client, error) {
	ethClient, err := ethclient.Dial(config.DialURI)
	if err != nil {
		return common.Address{}, &Client{}, err
	}

	privateKey, err := crypto.HexToECDSA(config.PrivateKey)
	if err != nil {
		return common.Address{}, &Client{}, err
	}

	chainID, err := ethClient.ChainID(ctx)
	if err != nil {
		return common.Address{}, &Client{}, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return common.Address{}, &Client{}, err
	}

	address, tx, registry, err := contracts.DeployNodeRegistry(auth, ethClient)
	if err != nil {
		return common.Address{}, &Client{}, err
	}

	client := &Client{
		Registry: registry,
		Auth:     auth,
		client:   ethClient,
		ticker:   time.NewTicker(200 * time.Millisecond),
	}

	return address, client, client.WaitForTx(ctx, tx.Hash())
}

// GetRelayer retrieves the current relayer address from the registry.
func (c *Client) GetRelayer() (string, [][]byte, error) {
	resp, err := c.Registry.GetRelayer(&bind.CallOpts{})
	if err != nil {
		return "", nil, err
	}

	return resp.Ip, resp.PublicKeys, nil
}

// GetResolver fetches the resolver address associated with the given public key.
func (c *Client) GetResolver(publicKey []byte) (string, error) {
	return c.Registry.GetResolver(&bind.CallOpts{}, publicKey)
}

// RegisterRelayer registers a new relayer with the specified IP address.
func (c *Client) RegisterRelayer(ctx context.Context, ipAddress string) error {
	tx, err := c.Registry.RegisterRelayer(c.Auth, ipAddress)
	if err != nil {
		return err
	}

	return c.WaitForTx(ctx, tx.Hash())
}

// RegisterResolver registers a new resolver with the given IP address and public key.
func (c *Client) RegisterResolver(ctx context.Context, ipAddress string, publicKey []byte) error {
	tx, err := c.Registry.RegisterResolver(c.Auth, ipAddress, publicKey)
	if err != nil {
		return err
	}

	return c.WaitForTx(ctx, tx.Hash())
}

// Close closes ethereum client.
func (c *Client) Close() {
	c.client.Close()
}

// WaitForTx block execution until transaction receipt is received or context is cancelled.
func (c *Client) WaitForTx(ctx context.Context, hash common.Hash) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return ErrContextCancelled
		case <-c.ticker.C:
			receipt, err := c.client.TransactionReceipt(ctx, hash)
			if err == nil {
				if receipt.Status == 1 {
					return nil
				}

				return ErrTransactionFailed
			}
			if !errors.Is(err, ethereum.NotFound) {
				return err
			}
		}
	}
}
