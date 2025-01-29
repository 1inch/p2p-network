package webrtc

import "time"

// RetryConfig represents the configuration for retry request to resolver
type RetryConfig struct {
	RequestCount    uint8
	RequestInterval time.Duration
}
