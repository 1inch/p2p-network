package webrtc

import "time"

// Retry represents the configuration for retry request to resolver
type Retry struct {
	// Count represents count of requests if resolver returned error
	Count uint8
	// Interval represents interval between requests
	Interval time.Duration
}

// PeerRangePort represents the configuration for peer connections port range between min and max
type PeerRangePort struct {
	Min uint16
	Max uint16
}
