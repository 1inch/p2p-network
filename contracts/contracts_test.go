//go:build deploy
// +build deploy

package contracts_test

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	registry "github.com/1inch/p2p-network/internal/registry"
)

const (
	rpcURL        = "http://127.0.0.1:8545"
	privateKeyHex = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	chainIDVal    = 31337
	contractAddr  = "0x5fbdb2315678afecb367f032d93f642f64180aa3"
)

func TestDeployContract(t *testing.T) {
	ctx := context.Background()
	config := &registry.Config{
		DialURI:    rpcURL,
		PrivateKey: privateKeyHex,
	}

	_, err := registry.DeployNodeRegistry(ctx, config)
	require.NoError(t, err, "contract deployment failed")
}

func TestRegisterResolver(t *testing.T) {
	ctx := context.Background()
	resolverIP := "http://127.0.0.1:8081"
	rpcURL := "http://127.0.0.1:8545"
	resolverPrivateKey := "5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a"

	privKey, err := crypto.HexToECDSA(resolverPrivateKey)
	require.NoError(t, err, "invalid private key")
	resolverPublicKeyBytes := crypto.CompressPubkey(&privKey.PublicKey)

	client, err := registry.Dial(ctx, &registry.Config{
		DialURI:         rpcURL,
		PrivateKey:      privateKeyHex,
		ContractAddress: contractAddr,
	})
	require.NoError(t, err, "failed to connect to %s", rpcURL)

	err = client.RegisterResolver(ctx, resolverIP, resolverPublicKeyBytes)
	require.NoError(t, err, "contract deployment failed")

	ip, err := client.Registry.GetResolver(&bind.CallOpts{}, resolverPublicKeyBytes)
	require.NoError(t, err, "contract get relayer failed")

	require.Equal(t, resolverIP, ip)
}
