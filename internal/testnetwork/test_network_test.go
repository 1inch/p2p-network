package testnetwork_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/1inch/p2p-network/internal/testnetwork"
)

const (
	awaitTimeout = 4 * time.Second
	retryBackoff = 500 * time.Millisecond
)

func TestNetworkWithOneNode(t *testing.T) {
	testnetwork.Run(t, 1, 1, func(tn *testnetwork.TestNetwork) {
		require.NotNil(t, tn)
		require.Equal(t, 1, tn.RelayerCount)
		require.Equal(t, 1, tn.ResolverCount)
	})
}

func TestNetworkWithMultipleNode(t *testing.T) {
	testnetwork.Run(t, 3, 3, func(tn *testnetwork.TestNetwork) {
		require.NotNil(t, tn)
		require.Equal(t, 3, tn.RelayerCount)
		require.Equal(t, 3, tn.RelayerCount)
	})
}

func TestNetworkStartStop(t *testing.T) {
	testnetwork.Run(t, 3, 3, func(tn *testnetwork.TestNetwork) {
		tn.Stop()

		for i := range tn.RelayerNodes {
			require.EventuallyWithT(t, func(collect *assert.CollectT) {
				assert.False(collect, testnetwork.IsPortBusy(tn.HTTPPorts[i]), "http port %d is still busy", tn.HTTPPorts[i])
			}, awaitTimeout, retryBackoff)
		}

		for i := range tn.ResolverNodes {
			require.EventuallyWithT(t, func(collect *assert.CollectT) {
				assert.False(collect, testnetwork.IsPortBusy(tn.GRPCPorts[i]), "grpc port %d is still busy", tn.GRPCPorts[i])
			}, awaitTimeout, retryBackoff)
		}
	})
}
