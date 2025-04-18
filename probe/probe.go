
// Package probe contains the probe implementation.
package probe

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/megaease/easeprobe/global"
)

// Prober Interface
type Prober interface {
	LabelMap() prometheus.Labels
	SetLabelMap(labels prometheus.Labels)
	Kind() string
	Name() string
	Channels() []string
	Timeout() time.Duration
	Interval() time.Duration
	Result() *Result
	Config(global.ProbeSettings) error
	Probe() Result
}
