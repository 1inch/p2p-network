package relayer

import "time"

// Config represents the configuration for the relayer node.
type Config struct {
	LogLevel             string             `yaml:"log_level"`
	HTTPEndpoint         string             `yaml:"http_endpoint"`
	WebRTCICEServer      string             `yaml:"webrtc_ice_server"`
	GRPCServerAddress    string             `yaml:"grpc_server_address"`
	WithNodeRegistry     bool               `yaml:"with_node_registry"`
	BlockchainRPCAddress string             `yaml:"blockchain_rpc_address"`
	ContractAddress      string             `yaml:"contract_address"`
	PrivateKey           string             `yaml:"private_key"`
	RetryRequestConfig   RetryRequestConfig `yaml:"retry_request"`
}

// RetryRequestConfig represents the configuration for retry request to resolver
type RetryRequestConfig struct {
	// Enabled represents retry is enabled/disabled
	Enabled bool `yaml:"enabled"`
	// Count represents count of requests if resolver returned error
	Count uint8 `yaml:"count"`
	// Interval represents interval between requests
	Interval time.Duration `yaml:"interval"`
}

// DefaultConfig returns default configuration for the relayer node.
func DefaultConfig() Config {
	return Config{
		LogLevel:             "DEBUG",
		HTTPEndpoint:         "127.0.0.1:0",
		WebRTCICEServer:      "stun:stun1.l.google.com:19302",
		GRPCServerAddress:    "127.0.0.1:0",
		WithNodeRegistry:     false,
		BlockchainRPCAddress: "http://127.0.0.1:8545",
		ContractAddress:      "0x8464135c8F25Da09e49BC8782676a84730C318bC",
		PrivateKey:           "59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d",
		RetryRequestConfig: RetryRequestConfig{
			Enabled: false,
		},
	}
}
