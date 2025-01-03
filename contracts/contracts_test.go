//go:build deploy
// +build deploy

package contracts_test

import (
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"

	registry "github.com/1inch/p2p-network/contracts"
)

const (
	rpcURL               = "http://127.0.0.1:8545"
	privateKeyHex        = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	chainIDVal           = 31337
	expectedContractAddr = "0x5fbdb2315678afecb367f032d93f642f64180aa3"
)

func TestDeployContract(t *testing.T) {
	client, err := ethclient.Dial(rpcURL)
	require.NoError(t, err, "failed to connect to %s", rpcURL)

	privKey, err := crypto.HexToECDSA(privateKeyHex)
	require.NoError(t, err, "invalid private key")

	auth, err := bind.NewKeyedTransactorWithChainID(privKey, big.NewInt(chainIDVal))
	require.NoError(t, err, "failed to create transactor")

	addr, tx, _, err := registry.DeployNodeRegistry(auth, client)
	require.NoError(t, err, "contract deployment failed")

	t.Logf("Deployed contract at: %s", addr.Hex())
	t.Logf("Deployment TX: %s", tx.Hash().Hex())

	require.Equal(t, expectedContractAddr, strings.ToLower(addr.Hex()),
		"deployed contract address should match the deterministic address")
}
