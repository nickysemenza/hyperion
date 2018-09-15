package metrics

import (
	"github.com/nickysemenza/hyperion/core/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

//Holds metric registrations
var (
	CueBacklogCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "hyperion",
		Name:      "cue_backlog_count",
		Help:      "Cue Size",
	}, []string{"cuestack_name"})

	CueProcessedCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "hyperion",
		Name:      "cue_processed_count",
		Help:      "Cue Size (Processed)",
	}, []string{"cuestack_name"})

	CueExecutionDrift = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "hyperion",
		Name:      "cue_execution_drift",
		Help:      "Drift in seconds of cue eecution",
	})

	ExternalResponseTime = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "hyperion",
		Name:      "external_response_time",
		Help:      "Speed of making requests to external sources (s)",
	}, []string{"target"})

	serverVersion = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "hyperion",
		Name:      "server_version",
		Help:      "currently runnning version",
	}, []string{"version"})
)

// Register inits prometheus registration
func init() {
	serverVersion.WithLabelValues(config.GetVersion()).Set(1)
}
