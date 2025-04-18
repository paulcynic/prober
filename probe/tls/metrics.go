
package tls

import (
	"crypto/tls"
	"time"

	"github.com/megaease/easeprobe/global"
	"github.com/megaease/easeprobe/metric"
	"github.com/prometheus/client_golang/prometheus"
)

// code and metric idea from https://github.com/prometheus/blackbox_exporter/blob/master/prober/tls.go
type metrics struct {
	EarliestCertExpiry              *prometheus.GaugeVec
	LastChainExpiryTimestampSeconds *prometheus.GaugeVec
}

// newMetrics create the HTTP metrics
func newMetrics(subsystem, name string, constLabels prometheus.Labels) *metrics {
	namespace := global.GetEaseProbe().Name
	return &metrics{
		EarliestCertExpiry: metric.NewGauge(namespace, subsystem, name, "earliest_cert_expiry",
			"last TLS chain expiry in timestamp seconds", []string{"endpoint"}, constLabels),
		LastChainExpiryTimestampSeconds: metric.NewGauge(namespace, subsystem, name, "last_chain_expiry_timestamp_seconds",
			"earliest TLS cert expiry in unix time", []string{"endpoint"}, constLabels),
	}
}

func getEarliestCertExpiry(state *tls.ConnectionState) time.Time {
	earliest := time.Time{}
	for _, cert := range state.PeerCertificates {
		if (earliest.IsZero() || cert.NotAfter.Before(earliest)) && !cert.NotAfter.IsZero() {
			earliest = cert.NotAfter
		}
	}
	return earliest
}

func getLastChainExpiry(state *tls.ConnectionState) time.Time {
	lastChainExpiry := time.Time{}
	for _, chain := range state.VerifiedChains {
		earliestCertExpiry := time.Time{}
		for _, cert := range chain {
			if (earliestCertExpiry.IsZero() || cert.NotAfter.Before(earliestCertExpiry)) && !cert.NotAfter.IsZero() {
				earliestCertExpiry = cert.NotAfter
			}
		}
		if lastChainExpiry.IsZero() || lastChainExpiry.Before(earliestCertExpiry) {
			lastChainExpiry = earliestCertExpiry
		}

	}
	return lastChainExpiry
}
