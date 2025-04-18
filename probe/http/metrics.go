
package http

import (
	"github.com/megaease/easeprobe/global"
	"github.com/megaease/easeprobe/metric"
	"github.com/prometheus/client_golang/prometheus"
)

// metrics is the metrics for http probe
type metrics struct {
	StatusCode       *prometheus.CounterVec
	ContentLen       *prometheus.GaugeVec
	DNSDuration      *prometheus.GaugeVec
	ConnectDuration  *prometheus.GaugeVec
	TLSDuration      *prometheus.GaugeVec
	SendDuration     *prometheus.GaugeVec
	WaitDuration     *prometheus.GaugeVec
	TransferDuration *prometheus.GaugeVec
	TotalDuration    *prometheus.GaugeVec
}

// newMetrics create the HTTP metrics
func newMetrics(subsystem, name string, constLabels prometheus.Labels) *metrics {
	namespace := global.GetEaseProbe().Name
	return &metrics{
		StatusCode: metric.NewCounter(namespace, subsystem, name, "status_code",
			"HTTP Status Code", []string{"name", "status", "endpoint"}, constLabels),
		ContentLen: metric.NewGauge(namespace, subsystem, name, "content_len",
			"HTTP Content Length", []string{"name", "status", "endpoint"}, constLabels),
		DNSDuration: metric.NewGauge(namespace, subsystem, name, "dns_duration",
			"DNS Duration", []string{"name", "status", "endpoint"}, constLabels),
		ConnectDuration: metric.NewGauge(namespace, subsystem, name, "connect_duration",
			"TCP Connection Duration", []string{"name", "status", "endpoint"}, constLabels),
		TLSDuration: metric.NewGauge(namespace, subsystem, name, "tls_duration",
			"TLS Duration", []string{"name", "status", "endpoint"}, constLabels),
		SendDuration: metric.NewGauge(namespace, subsystem, name, "send_duration",
			"Send Duration", []string{"name", "status", "endpoint"}, constLabels),
		WaitDuration: metric.NewGauge(namespace, subsystem, name, "wait_duration",
			"Wait Duration", []string{"name", "status", "endpoint"}, constLabels),
		TransferDuration: metric.NewGauge(namespace, subsystem, name, "transfer_duration",
			"Transfer Duration", []string{"name", "status", "endpoint"}, constLabels),
		TotalDuration: metric.NewGauge(namespace, subsystem, name, "total_duration",
			"Total Duration", []string{"name", "status", "endpoint"}, constLabels),
	}
}
