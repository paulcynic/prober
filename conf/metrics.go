package conf

import (
	"github.com/megaease/easeprobe/global"
	"github.com/megaease/easeprobe/metric"
	"github.com/prometheus/client_golang/prometheus"
)

// ConfigMetrics is the conf metrics
type ConfigMetrics struct {
	Timestamp *prometheus.GaugeVec
	Status    *prometheus.GaugeVec
}

// NewConfigMetrics create the metrics
func NewConfigMetrics() *ConfigMetrics {
	namespace := global.DefaultMetricNamespace
	subsystem := "config"
	return &ConfigMetrics{
		Timestamp: metric.NewGauge(
			namespace,
			subsystem,
			"last_modified",
			"timestamp",
			"Unix timestamp (in seconds) when the config was last successfully fetched via HTTP",
			[]string{"endpoint"},
			prometheus.Labels{},
		),
		Status: metric.NewGauge(
			namespace,
			subsystem,
			"availability",
			"status",
			"Whether the config is currently available via HTTP (1=available, 0=not available)",
			[]string{"endpoint"},
			prometheus.Labels{},
		),
	}
}
