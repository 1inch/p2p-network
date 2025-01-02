package contracts

import (
	"context"
	"errors"
	"time"

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

// Client represents storage client.
type Client struct {
	Registry *NodeRegistry
	Auth     *bind.TransactOpts
	client   *ethclient.Client
	ticker   *time.Ticker
}

// Dial creates eth client, new smart-contract instance, auth.
func Dial(ctx context.Context, url, key, contractAddress string) (*Client, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return &Client{}, err
	}

	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		return &Client{}, err
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return &Client{}, err
	}

	registry, err := NewNodeRegistry(common.HexToAddress(contractAddress), client)
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
