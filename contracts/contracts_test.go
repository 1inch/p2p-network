//go:build deploy
// +build deploy

package contracts_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	registry "github.com/1inch/p2p-network/internal/registry"
)

const (
	rpcURL        = "http://127.0.0.1:8545"
	privateKeyHex = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
)

func TestDeployContract(t *testing.T) {
	ctx := context.Background()
	config := registry.Config{
		DialURI:    rpcURL,
		PrivateKey: privateKeyHex,
	}

	_, err := registry.DeployNodeRegistry(ctx, config)
	require.NoError(t, err, "contract deployment failed")
}
