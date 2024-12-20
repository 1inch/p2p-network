package relayer

// Config represents the configuration for the relayer node.
type Config struct {
	LogLevel          string `yaml:"log_level"`
	HTTPEndpoint      string `yaml:"http_endpoint"`
	WebRTCICEServer   string `yaml:"webrtc_ice_server"`
	GRPCServerAddress string `yaml:"grpc_server_address"`
}

// DefaultConfig returns default configuration for the relayer node.
func DefaultConfig() Config {
	return Config{
		LogLevel:          "DEBUG",
		HTTPEndpoint:      "127.0.0.1:0",
		WebRTCICEServer:   "stun:stun1.l.google.com:19302",
		GRPCServerAddress: "127.0.0.1:0",
	}
}
