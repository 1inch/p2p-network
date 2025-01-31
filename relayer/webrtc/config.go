package webrtc

import "time"

// Config represents the configuration for webrtc server
type Config struct {
	RetryRequestConfig RetryRequestConfig
	PortRangeConfig    PortRangeConfig
}

// RetryRequestConfig represents the configuration for retry request to resolver
type RetryRequestConfig struct {
	Enabled  bool
	Count    uint8
	Interval time.Duration
}

// PortRangeConfig represents the configuration for peer connections port range between min and max
type PortRangeConfig struct {
	Enabled bool
	Min     uint16
	Max     uint16
}
