package conf

import (
	"github.com/megaease/easeprobe/global"
	"github.com/megaease/easeprobe/metric"
	"github.com/prometheus/client_golang/prometheus"
)

// ConfigMetrics is the conf metrics
type ConfigMetrics struct {
	Duration *prometheus.GaugeVec
	Status   *prometheus.GaugeVec
}

// NewConfigMetrics create the metrics
func NewConfigMetrics(constLabels prometheus.Labels) *ConfigMetrics {
	namespace := global.GetEaseProbe().Name
	subsystem := "config"
	name := "availability"
	return &ConfigMetrics{
		Duration: metric.NewGauge(
			namespace,
			subsystem,
			name,
			"timestamp",
			"Unix timestamp (in seconds) when the config was last successfully fetched via HTTP",
			[]string{"endpoint"},
			constLabels,
		),
		Status: metric.NewGauge(
			namespace,
			subsystem,
			name,
			"status",
			"Whether the config is currently available via HTTP (1=available, 0=not available)",
			[]string{"endpoint"},
			constLabels,
		),
	}
}
