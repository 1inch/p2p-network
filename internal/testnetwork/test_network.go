// Package testnetwork provides a test network for running multiple nodes.
package testnetwork

import (
	"context"
	"log/slog"
	"net"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/1inch/p2p-network/internal/registry"
	"github.com/1inch/p2p-network/relayer"
	"github.com/1inch/p2p-network/resolver"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	defaultAwaitTimeout = 4 * time.Second
	defaultRetryBackoff = 500 * time.Millisecond
)

var (
	relayerPrivateKeys = []string{
		"59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d",
		"5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a",
		"7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6",
		"47e179ec197488593b187f80a00eb0da91f1b9d0b13f8733639f19c30a34926a",
		"8b3a350cf5c34c9194ca85829a2df0ec3153be0318b5e2d3348e872092edffba",
		"92db14e403b83dfe3df233f83dfa3a0d7096f21ca9b0d6d6b8d88b2b4ec1564e",
		"4bbbf85ce3377467afe5d46f804f221813b2bb87f24d81f60f1fcdbf7cbf4356",
		"dbda1821b80551c9d65939329250298aa3472ba22feea921c0cf5d620ea67b97",
		"2a871d0798f97d79848a013d4936a73bf4cc922c825d33c1cf7073dff6d409c6",
		"f214f2b2cd398c806f84e317254e0f0b801d0643303237d97a22a48e01628897",
	}

	resolverPrivateKeys = []string{
		"701b615bbdfb9de65240bc28bd21bbc0d996645a3dd57e7b12bc2bdf6f192c82",
		"a267530f49f8280200edf313ee7af6b827f2a8bce2897751d06a843f644967b1",
		"47c99abed3324a2707c28affff1267e45918ec8c3f20b8aa892e8b065d2942dd",
		"c526ee95bf44d8fc405a158bb884d9d1238d99f0612e9f33d006bb0789009aaa",
		"8166f546bab6da521a8369cab06c5d2b9e46670292d85c875ee9ec20e84ffb61",
		"ea6c44ac03bff858b476bba40716402b03e41b8e97e276d1baec7c37d42484a0",
		"689af8efa8c651a91ad287602527f3af2fe9f6501a7ac4b061667b5a93e037fd",
		"de9be858da4a475276426320d5e9262ecfc3ba460bfac56360bfa6c4c28b4ee0",
		"df57089febbacf7ba0bc227dafbffa9fc08a93fdc68e1e42411a14efcf23656e",
		"eaa861a9a01391ed3d587d8a5a84ca56ee277629a8b02c22093a419bf240e65d",
	}
)

// Config is the configuration for TestNetwork.
type Config struct {
	WithNodeRegistry   bool
	ResolverApiConfigs resolver.ApiConfigs
}

// Option some option for TestNetwork.
type Option func(config *Config)

// TestNetwork is a test network for running multiple nodes.
type TestNetwork struct {
	t                   *testing.T
	cancelFns           context.CancelFunc
	wg                  sync.WaitGroup
	RelayerCount        int
	ResolverCount       int
	ResolverNodes       []*resolver.Resolver
	RelayerNodes        []*relayer.Relayer
	GRPCPorts           []int
	HTTPPorts           []int
	ResolverPrivateKeys []string
	RelayerPrivateKeys  []string
}

// WithInfura option to set the infura apis.
func WithInfura(key string) Option {
	return func(cfg *Config) {
		cfg.ResolverApiConfigs.Infura.Enabled = true
		cfg.ResolverApiConfigs.Infura.Key = key
	}
}

// WithNodeRegistry option to enable the node registry.
func WithNodeRegistry() Option {
	return func(cfg *Config) {
		cfg.WithNodeRegistry = true
	}
}

// New creates a new test network.
func New(ctx context.Context, t *testing.T, relayerCount, resolverCount int, options ...Option) *TestNetwork {
	testNetwork := &TestNetwork{
		t:                   t,
		RelayerCount:        relayerCount,
		ResolverCount:       resolverCount,
		RelayerNodes:        make([]*relayer.Relayer, relayerCount),
		ResolverNodes:       make([]*resolver.Resolver, resolverCount),
		GRPCPorts:           make([]int, resolverCount),
		HTTPPorts:           make([]int, relayerCount),
		ResolverPrivateKeys: resolverPrivateKeys,
		RelayerPrivateKeys:  relayerPrivateKeys,
	}

	tnCfg := &Config{}

	for _, opt := range options {
		opt(tnCfg)
	}

	cfg := relayer.DefaultConfig()

	if tnCfg.WithNodeRegistry {
		cfg.DiscoveryConfig.WithNodeRegistry = tnCfg.WithNodeRegistry
	}

	_, _, err := registry.DeployNodeRegistry(context.Background(), &registry.Config{
		DialURI:         cfg.DiscoveryConfig.RpcUrl,
		PrivateKey:      cfg.PrivateKey,
		ContractAddress: cfg.DiscoveryConfig.ContractAddress,
	})
	assert.NoError(t, err, "Failed to create registry client")

	var level = new(slog.LevelVar)
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	}))
	level.Set(slog.LevelDebug)

	for i := 0; i < relayerCount; i++ {
		relayerNode, err := relayer.New(&cfg, logger)
		require.NoError(testNetwork.t, err)

		port, err := parsePort(relayerNode.HTTPServer.Addr())
		require.NoError(testNetwork.t, err)

		testNetwork.HTTPPorts[i] = port

		testNetwork.RelayerNodes[i] = relayerNode
	}

	for i := 0; i < resolverCount; i++ {
		resolverCfg := getResolverConfig()
		resolverCfg.PrivateKey = resolverPrivateKeys[i]

		if tnCfg.ResolverApiConfigs.Infura.Enabled {
			resolverCfg.Apis.Infura.Enabled = tnCfg.ResolverApiConfigs.Infura.Enabled
			resolverCfg.Apis.Infura.Key = tnCfg.ResolverApiConfigs.Infura.Key
		}

		resolverNode, err := resolver.New(resolverCfg, logger)
		require.NoError(testNetwork.t, err)

		port, err := parsePort(resolverNode.Addr())
		require.NoError(testNetwork.t, err)

		testNetwork.GRPCPorts[i] = port
		registerResolver(ctx, t, i, cfg, resolverNode.Addr())

		testNetwork.ResolverNodes[i] = resolverNode
	}

	return testNetwork
}

// Start starts the test network.
func (tn *TestNetwork) Start(ctx context.Context) {
	require.True(tn.t, tn.RelayerCount > 0, "relayer node count must be greater than 0")
	require.True(tn.t, tn.ResolverCount > 0, "resolver node count must be greater than 0")

	ctx, cancel := context.WithCancel(context.Background())
	tn.cancelFns = cancel

	for _, node := range tn.RelayerNodes {
		tn.wg.Add(1)

		go func(node *relayer.Relayer) {
			defer tn.wg.Done()

			err := node.Run(ctx)
			require.NoError(tn.t, err)
		}(node)
	}

	for _, node := range tn.ResolverNodes {
		tn.wg.Add(1)

		go func(node *resolver.Resolver) {
			defer tn.wg.Done()

			err := node.Run()
			require.NoError(tn.t, err)

			<-ctx.Done()
			err = node.Stop()
			require.NoError(tn.t, err)
		}(node)
	}
}

// Stop stops the test network.
func (tn *TestNetwork) Stop() {
	tn.cancelFns()
	tn.wg.Wait()
}

// Run starts a test network and runs the given function.
func Run(t *testing.T, relayerCount, resolverCount int, useFunc func(network *TestNetwork), options ...Option) {
	ctx := context.Background()
	tn := New(ctx, t, relayerCount, resolverCount, options...)

	tn.Start(ctx)
	defer tn.Stop()

	for i := range tn.RelayerNodes {
		tn.waitRelayer(i)
	}

	for i := range tn.ResolverNodes {
		tn.waitResolver(i)
	}

	useFunc(tn)
}

func (tn *TestNetwork) waitRelayer(nodeIndex int) {
	require.EventuallyWithT(tn.t, func(collect *assert.CollectT) {
		assert.True(collect, IsPortBusy(tn.HTTPPorts[nodeIndex]), "http port %d is not busy", tn.HTTPPorts[nodeIndex])
	}, defaultAwaitTimeout, defaultRetryBackoff)
}

func (tn *TestNetwork) waitResolver(nodeIndex int) {
	require.EventuallyWithT(tn.t, func(collect *assert.CollectT) {
		assert.True(collect, IsPortBusy(tn.GRPCPorts[nodeIndex]), "grpc port %d is not busy", tn.GRPCPorts[nodeIndex])
	}, defaultAwaitTimeout, defaultRetryBackoff)
}

// IsPortBusy checks if the port is busy.
func IsPortBusy(port int) bool {
	conn, err := net.Dial("tcp", "127.0.0.1:"+strconv.Itoa(port))
	if err != nil {
		return false
	}

	_, _ = conn.Write([]byte{})
	_ = conn.Close()
	return true
}

func parsePort(addr string) (int, error) {
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		return 0, err
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		return 0, err
	}

	return portInt, nil
}

func getResolverConfig() resolver.Config {
	return resolver.Config{
		GrpcEndpoint: "localhost:0",
		LogLevel:     slog.LevelInfo,
		Apis: resolver.ApiConfigs{
			Default: resolver.DefaultApiConfig{
				Enabled: true,
			},
		},
	}
}

func registerResolver(ctx context.Context, t *testing.T, index int, cfg relayer.Config, ipAddress string) {
	privateKey := resolverPrivateKeys[index]
	privKey, err := crypto.HexToECDSA(privateKey)
	require.NoError(t, err, "invalid private key")
	publicKey := crypto.CompressPubkey(&privKey.PublicKey)

	client, err := registry.Dial(ctx, &registry.Config{
		DialURI:         cfg.DiscoveryConfig.RpcUrl,
		PrivateKey:      privateKey,
		ContractAddress: cfg.DiscoveryConfig.ContractAddress,
	})
	require.NoError(t, err, "failed to connect to %s", cfg.DiscoveryConfig.RpcUrl)

	_ = client.RegisterResolver(ctx, ipAddress, publicKey)
}
