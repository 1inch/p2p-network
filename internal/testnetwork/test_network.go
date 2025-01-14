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

	"github.com/1inch/p2p-network/relayer"
	"github.com/1inch/p2p-network/resolver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

const (
	defaultAwaitTimeout = 4 * time.Second
	defaultRetryBackoff = 500 * time.Millisecond
)

type Config struct {
	WithNodeRegistry   bool
	ResolverApiConfigs resolver.ApiConfigs
}

// Option some option for TestNetwork.
type Option func(config *Config)

// TestNetwork is a test network for running multiple nodes.
type TestNetwork struct {
	t             *testing.T
	cancelFns     context.CancelFunc
	wg            sync.WaitGroup
	RelayerCount  int
	ResolverCount int
	ResolverNodes []*grpc.Server
	RelayerNodes  []*relayer.Relayer
	GRPCPorts     []int
	HTTPPorts     []int
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
		t:             t,
		RelayerCount:  relayerCount,
		ResolverCount: resolverCount,
		RelayerNodes:  make([]*relayer.Relayer, relayerCount),
		ResolverNodes: make([]*grpc.Server, resolverCount),
		GRPCPorts:     make([]int, resolverCount),
		HTTPPorts:     make([]int, relayerCount),
	}

	tnCfg := &Config{}

	for _, opt := range options {
		opt(tnCfg)
	}

	for i := 0; i < relayerCount; i++ {
		cfg := relayer.DefaultConfig()

		if tnCfg.WithNodeRegistry {
			cfg.WithNodeRegistry = tnCfg.WithNodeRegistry
		}

		logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

		relayerNode, err := relayer.New(&cfg, logger)
		require.NoError(testNetwork.t, err)

		port, err := parsePort(relayerNode.HTTPServer.Addr())
		require.NoError(testNetwork.t, err)

		testNetwork.HTTPPorts[i] = port

		testNetwork.RelayerNodes[i] = relayerNode
	}

	for i := 0; i < resolverCount; i++ {
		cfg := getResolverConfig()

		if tnCfg.ResolverApiConfigs.Infura.Enabled {
			cfg.Apis.Infura.Enabled = tnCfg.ResolverApiConfigs.Infura.Enabled
			cfg.Apis.Infura.Key = tnCfg.ResolverApiConfigs.Infura.Key
		}

		resolverNode, address, err := resolver.Run(cfg)
		require.NoError(testNetwork.t, err)

		port, err := parsePort(address)
		require.NoError(testNetwork.t, err)

		testNetwork.GRPCPorts[i] = port

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

	for i, node := range tn.RelayerNodes {
		tn.wg.Add(1)

		go func(i int, node *relayer.Relayer) {
			defer tn.wg.Done()

			err := node.Run(ctx)
			require.NoError(tn.t, err)
		}(i, node)
	}

	for i, node := range tn.ResolverNodes {
		tn.wg.Add(1)

		go func(i int, node *grpc.Server) {
			defer tn.wg.Done()

			<-ctx.Done()
			node.GracefulStop()
		}(i, node)
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

func getResolverConfig() *resolver.Config {
	return &resolver.Config{
		Port:     0,
		LogLevel: slog.LevelInfo,
		Apis: resolver.ApiConfigs{
			Default: resolver.DefaultApiConfig{
				Enabled: true,
			},
		},
	}
}
