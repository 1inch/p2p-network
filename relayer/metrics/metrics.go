// metrics.go
package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// HTTP instrumentation
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "relayer_http_requests_total",
			Help: "Total number of HTTP requests received",
		},
		[]string{"handler", "method", "status"},
	)
	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "relayer_http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: []float64{0.05, 0.5, 1, 2, 3, 5, 10, 15, 30, 60},
		},
		[]string{"handler", "method"},
	)

	// ICE candidate instrumentation
	IceCandidateSentTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "relayer_ice_candidate_sent_total",
			Help: "Total number of ICE candidates sent",
		},
		[]string{"session_id", "status"},
	)
	IceCandidateSendDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "relayer_ice_candidate_send_duration_seconds",
			Help:    "Duration for sending ICE candidates in seconds",
			Buckets: []float64{0.05, 0.5, 1, 2, 3, 5, 10, 15, 30, 60},
		},
		[]string{"session_id"},
	)

	// Active PeerConnections
	ActivePeerConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "relayer_active_peer_connections",
			Help: "Current number of active PeerConnections",
		},
	)

	// SDP negotiation metrics
	SdpNegotiationTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "relayer_sdp_negotiation_total",
			Help: "Total number of SDP negotiations",
		},
		[]string{"status"},
	)
	SdpNegotiationDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "relayer_sdp_negotiation_duration_seconds",
			Help:    "Duration of SDP negotiations in seconds",
			Buckets: []float64{0.05, 0.5, 1, 2, 3, 5, 10, 15, 30, 60},
		},
	)

	// gRPC call metrics
	GrpcRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "relayer_grpc_requests_total",
			Help: "Total number of gRPC requests",
		},
		[]string{"method", "status"},
	)
	GrpcRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "relayer_grpc_request_duration_seconds",
			Help:    "Duration of gRPC requests in seconds",
			Buckets: []float64{0.05, 0.5, 1, 2, 3, 5, 10, 15, 30, 60},
		},
		[]string{"method"},
	)

	// Data channel metrics
	DataChannelMessagesSent = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "relayer_data_channel_messages_sent_total",
			Help: "Total number of messages sent over data channels",
		},
		[]string{"session_id", "status"},
	)
	DataChannelMessagesReceived = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "relayer_data_channel_messages_received_total",
			Help: "Total number of messages received over data channels",
		},
		[]string{"session_id"},
	)
	DataChannelLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "relayer_data_channel_latency_seconds",
			Help:    "Time taken to process data channel messages in seconds",
			Buckets: []float64{0.05, 0.5, 1, 2, 3, 5, 10, 15, 30, 60},
		},
		[]string{"session_id"},
	)

	// End-to-end workflow metrics
	EndToEndWorkflowLatency = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "relayer_end_to_end_workflow_latency_seconds",
			Help:    "Latency of end-to-end workflow from receiving SDP offer to sending answer, in seconds",
			Buckets: []float64{0.05, 0.5, 1, 2, 3, 5, 10, 15, 30, 60},
		},
	)
	EndToEndWorkflowCompleted = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "relayer_end_to_end_workflow_completed_total",
			Help: "Total number of completed end-to-end workflows",
		},
	)
)

func init() {
	prometheus.MustRegister(
		HttpRequestsTotal, HttpRequestDuration,
		IceCandidateSentTotal, IceCandidateSendDuration,
		ActivePeerConnections,
		SdpNegotiationTotal, SdpNegotiationDuration,
		GrpcRequestsTotal, GrpcRequestDuration,
		DataChannelMessagesSent, DataChannelMessagesReceived,
		DataChannelLatency, EndToEndWorkflowLatency,
		EndToEndWorkflowCompleted,
	)
}

func MetricsHandler() http.Handler {
	return promhttp.Handler()
}
