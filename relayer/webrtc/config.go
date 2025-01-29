package webrtc

import "time"

// RetryRequestConfig represents the configuration for retry request to resolver
type RetryRequestConfig struct {
	Enabled  bool
	Count    uint8
	Interval time.Duration
}
