package webrtc

import "time"

// Config represents the configuration for webrtc server
type Config struct {
	ICEServer      string
	RetryConfig    RetryConfig
	PeerPortConfig PeerPortConfig
}

// RetryConfig represents the configuration for retry request to resolver
type RetryConfig struct {
	// Enabled represents retry is enabled/disabled
	Enabled bool
	// Count represents count of requests if resolver returned error
	Count uint8
	// Interval represents interval between requests
	Interval time.Duration
}

// PortConfig represents the configuration for peer connections port range between min and max
type PeerPortConfig struct {
	Enabled bool
	Min     uint16
	Max     uint16
}
