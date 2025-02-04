package webrtc

import "time"

// RetryOption represents the configuration for retry request to resolver
type RetryOption struct {
	// Count represents count of requests if resolver returned error
	Count uint8
	// Interval represents interval between requests
	Interval time.Duration
}

// PeerPortOption represents the configuration for peer connections port range between min and max
type PeerPortOption struct {
	Min uint16
	Max uint16
}
