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
	contractAddr  = "0x5fbdb2315678afecb367f032d93f642f64180aa3"
)

func TestDeployContract(t *testing.T) {
	ctx := context.Background()
	config := &registry.Config{
		DialURI:    rpcURL,
		PrivateKey: privateKeyHex,
	}

	address, _, err := registry.DeployNodeRegistry(ctx, config)
	t.Logf("contract address: %s", address)
	require.NoError(t, err, "contract deployment failed")
}

func TestRegisterResolver(t *testing.T) {
	ctx := context.Background()
	resolverIP := "127.0.0.1:8001"
	relayerIP := "127.0.0.1:8880"
	rpcURL := "http://127.0.0.1:8545"
	resolverPrivateKey := "5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a"

	t.Logf("discovery mock contract address: %s", contractAddr)

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
	require.NoError(t, err, "resolver registration failed")

	ip, err := client.Registry.GetResolver(&bind.CallOpts{}, resolverPublicKeyBytes)
	require.NoError(t, err, "contract get resolver failed")

	require.Equal(t, resolverIP, ip)

	err = client.RegisterRelayer(ctx, relayerIP)
	require.NoError(t, err, "relayer registration failed")

	t.Log("resolver successfully registered")
}
