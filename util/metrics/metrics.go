package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

//Holds metric registrations
var (
	CueBacklogCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cue_backlog_count",
		Help: "Cue Size",
	}, []string{"cuestack_name"})

	CueProcessedCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cue_processed_count",
		Help: "Cue Size (Processed)",
	}, []string{"cuestack_name"})

	CueExecutionDriftNs = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cue_execution_drift_ns",
		Help: "Drift in ns of cue eecution",
	})

	ResponseTimeNsOLA = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "external_response_time_ola",
		Help: "Speed of making requests to ola (ns)",
	})

	ResponseTimeNsHue = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "external_response_time_hue",
		Help: "Speed of making requests to hue (ns)",
	})
)

// Register inits prometheus registration
func Register() {
	prometheus.MustRegister(CueBacklogCount)
	prometheus.MustRegister(CueProcessedCount)
	prometheus.MustRegister(CueExecutionDriftNs)
	prometheus.MustRegister(ResponseTimeNsOLA)
	prometheus.MustRegister(ResponseTimeNsHue)
}

//SetGagueWithNsFromTime is used for updating prometheus gague with time elapsed
func SetGagueWithNsFromTime(start time.Time, metric prometheus.Gauge) {
	metric.Set(float64(time.Since(start) / time.Nanosecond))
}
