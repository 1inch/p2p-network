package resolver

import "log/slog"

// DefaultApiConfig provides configuration for default api handler
type DefaultApiConfig struct {
	Enabled bool
}

// InfuraApiConfig provides configuration for Infura api handler
type InfuraApiConfig struct {
	Key     string
	Enabled bool
}

// ApiConfigs contains API-related configs
type ApiConfigs struct {
	Default DefaultApiConfig
	Infura  InfuraApiConfig
}

// Config represents resolver server config
type Config struct {
	// gRPC server port
	Port int

	// Can be one or more of the following: default,infura
	Apis ApiConfigs

	// Default loglevel
	LogLevel slog.Level
}
