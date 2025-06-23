package conf

import (
	"github.com/megaease/easeprobe/global"
	"github.com/megaease/easeprobe/metric"
	"github.com/prometheus/client_golang/prometheus"
)

// metrics is the conf metrics
type metrics struct {
	Duration *prometheus.GaugeVec
	Status   *prometheus.GaugeVec
}

// newMetrics create the metrics
func newMetrics(constLabels prometheus.Labels) *metrics {
	namespace := global.GetEaseProbe().Name
	subsystem := "config"
	name := "availability"
	return &metrics{
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
