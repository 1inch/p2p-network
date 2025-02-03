package relayer

import "time"

// Config represents the configuration for the relayer node.
type Config struct {
	LogLevel        string          `yaml:"log_level"`
	HTTPEndpoint    string          `yaml:"http_endpoint"`
	PrivateKey      string          `yaml:"private_key"`
	DiscoveryConfig DiscoveryConfig `yaml:"discovery"`
	WebrtcConfig    WebrtcConfig    `yaml:"webrtc"`
}

// WebrtcConfig represents the configuration for webrtc server
type WebrtcConfig struct {
	ICEServer      string         `yaml:"ice_server"`
	RetryConfig    RetryConfig    `yaml:"retry"`
	PeerPortConfig PeerPortConfig `yaml:"port"`
}

// DiscoveryConfig represents the configuration for discovery service
type DiscoveryConfig struct {
	RpcUrl           string `yaml:"rpc_url"`
	WithNodeRegistry bool   `yaml:"with_node_registry"`
	ContractAddress  string `yaml:"contract_address"`
}

// RetryConfig represents the configuration for retry request to resolver
type RetryConfig struct {
	// Enabled represents retry is enabled/disabled
	Enabled bool `yaml:"enabled"`
	// Count represents count of requests if resolver returned error
	Count uint8 `yaml:"count"`
	// Interval represents interval between requests
	Interval time.Duration `yaml:"interval"`
}

// PeerPortConfig represents the configuration for peer connections port range between min and max
type PeerPortConfig struct {
	Enabled bool   `yaml:"enabled"`
	Min     uint16 `yaml:"min"`
	Max     uint16 `yaml:"max"`
}

// DefaultConfig returns default configuration for the relayer node.
func DefaultConfig() Config {
	return Config{
		LogLevel:     "DEBUG",
		HTTPEndpoint: "127.0.0.1:0",
		PrivateKey:   "59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d",
		DiscoveryConfig: DiscoveryConfig{
			RpcUrl:           "http://127.0.0.1:8545",
			WithNodeRegistry: false,
			ContractAddress:  "0x5fbdb2315678afecb367f032d93f642f64180aa3",
		},
		WebrtcConfig: WebrtcConfig{
			ICEServer: "stun:stun1.l.google.com:19302",
			RetryConfig: RetryConfig{
				Enabled: false,
			},
			PeerPortConfig: PeerPortConfig{
				Enabled: false,
			},
		},
	}
}
