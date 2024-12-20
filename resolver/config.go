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

// Config represents resolver server config
type Config struct {
	// gRPC server port
	Port int `yaml:"port"`

	// Can be one or more of the following: default,infura
	Apis ApiConfigs `yaml:"apis"`

	// Default loglevel
	LogLevel slog.Level `yaml:"log_level"`
}
