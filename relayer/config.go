package relayer

// Config represents the configuration for the relayer node.
type Config struct {
	LogLevel             string `yaml:"log_level"`
	HTTPEndpoint         string `yaml:"http_endpoint"`
	WebRTCICEServer      string `yaml:"webrtc_ice_server"`
	GRPCServerAddress    string `yaml:"grpc_server_address"`
	BlockchainRPCAddress string `yaml:"blockchain_rpc_address"`
	ContractAddress      string `yaml:"contract_address"`
	PrivateKey           string `yaml:"private_key"`
	ChainID              int64  `yaml:"chain_id"`
}

// DefaultConfig returns default configuration for the relayer node.
func DefaultConfig() Config {
	return Config{
		LogLevel:             "DEBUG",
		HTTPEndpoint:         "127.0.0.1:0",
		WebRTCICEServer:      "stun:stun1.l.google.com:19302",
		GRPCServerAddress:    "127.0.0.1:0",
		BlockchainRPCAddress: "127.0.0.1:8545",
		ContractAddress:      "0x5fbdb2315678afecb367f032d93f642f64180aa3",
		PrivateKey:           "59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d",
		ChainID:              31337,
	}
}
