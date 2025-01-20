package resolver

import "log/slog"

// DefaultApiConfig provides configuration for default api handler
type DefaultApiConfig struct {
	Enabled bool `yaml:"enabled"`
}

// InfuraApiConfig provides configuration for Infura api handler
type InfuraApiConfig struct {
	Key     string `yaml:"key"`
	Enabled bool   `yaml:"enabled"`
}

// ApiConfigs contains API-related configs
type ApiConfigs struct {
	Default DefaultApiConfig `yaml:"default"`
	Infura  InfuraApiConfig  `yaml:"infura"`
}

// MetricConfig contain params for configure metrics
type MetricConfig struct {
	Enabled bool `yaml:"enabled"`
	Port    uint `yaml:"port"`
}

// Config represents resolver server config
type Config struct {

	// gRPC server endpoint
	GrpcEndpoint string `yaml:"grpc_endpoint"`

	// Default resolver node key
	PrivateKey string `yaml:"private_key"`

	// Discovery contract address
	ContractAddress string `yaml:"contract_address"`

	// rpc url to blockchain node
	RpcUrl string `yaml:"rpc_url"`

	// Can be one or more of the following: default,infura
	Apis ApiConfigs `yaml:"apis"`

	// Default loglevel
	LogLevel slog.Level `yaml:"log_level"`

	// Configuration metric
	Metric MetricConfig `yaml:"metric"`
}
