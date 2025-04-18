
package base

import (
	"github.com/megaease/easeprobe/global"
	"github.com/megaease/easeprobe/metric"
	"github.com/prometheus/client_golang/prometheus"
)

// metrics is the probe metrics
type metrics struct {
	TotalCnt  *prometheus.GaugeVec
	TotalTime *prometheus.GaugeVec
	Duration  *prometheus.GaugeVec
	Status    *prometheus.GaugeVec
	SLA       *prometheus.GaugeVec
}

// newMetrics create the metrics
func newMetrics(subsystem, name string, constLabels prometheus.Labels) *metrics {
	namespace := global.GetEaseProbe().Name
	return &metrics{
		TotalCnt: metric.NewGauge(namespace, subsystem, name, "total",
			"Total Probed Counts", []string{"name", "status", "endpoint"}, constLabels),
		TotalTime: metric.NewGauge(namespace, subsystem, name, "total_time",
			"Total Time(Seconds) of Status", []string{"name", "status", "endpoint"}, constLabels),
		Duration: metric.NewGauge(namespace, subsystem, name, "duration",
			"Probe Duration", []string{"name", "status", "endpoint"}, constLabels),
		Status: metric.NewGauge(namespace, subsystem, name, "status",
			"Probe Status", []string{"name", "endpoint"}, constLabels),
		SLA: metric.NewGauge(namespace, subsystem, name, "sla",
			"Probe SLA", []string{"name", "endpoint"}, constLabels),
	}
}
