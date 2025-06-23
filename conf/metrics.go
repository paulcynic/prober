package conf

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ConfigFetchFailure = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "easeprobe_config_fetch_failure_total",
			Help: "Total number of failed config fetch attempts via HTTP",
		},
		[]string{"url"},
	)
)

func init() {
	prometheus.MustRegister(ConfigFetchFailure)
}
