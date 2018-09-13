package metrics

import (
	"time"

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

	CueExecutionDriftNs = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "hyperion",
		Name:      "cue_execution_drift_ns",
		Help:      "Drift in ns of cue eecution",
	})

	ResponseTimeNsOLA = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "hyperion",
		Name:      "external_response_time_ola",
		Help:      "Speed of making requests to ola (ns)",
	})

	ResponseTimeNsHue = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "hyperion",
		Name:      "external_response_time_hue",
		Help:      "Speed of making requests to hue (ns)",
	})
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

//SetGagueWithNsFromTime is used for updating prometheus gague with time elapsed
func SetGagueWithNsFromTime(start time.Time, metric prometheus.Gauge) {
	metric.Set(float64(time.Since(start) / time.Nanosecond))
}
