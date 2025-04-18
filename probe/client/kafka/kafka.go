
// Package kafka is the native client probe for kafka.
package kafka

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/megaease/easeprobe/probe/client/conf"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	log "github.com/sirupsen/logrus"
)

// Kind is the type of driver
const Kind string = "Kafka"

// Kafka is the Kafka client
type Kafka struct {
	conf.Options `yaml:",inline"`
	tls          *tls.Config     `yaml:"-" json:"-"`
	Context      context.Context `yaml:"-" json:"-"`
}

// New create a Kafka client
func New(opt conf.Options) (*Kafka, error) {
	tls, err := opt.TLS.Config()
	if err != nil {
		log.Errorf("[%s / %s / %s] - TLS Config Error - %v", opt.ProbeKind, opt.ProbeName, opt.ProbeTag, err)
		return nil, fmt.Errorf("TLS Config Error - %v", err)
	}
	k := &Kafka{
		Options: opt,
		tls:     tls,
		Context: context.Background(),
	}
	return k, nil
}

// Kind return the name of client
func (k *Kafka) Kind() string {
	return Kind
}

// Probe do the health check
func (k *Kafka) Probe() (bool, string) {

	var dialer *kafka.Dialer

	if len(k.Password) > 0 {
		dialer = &kafka.Dialer{
			Timeout: k.Timeout(),
			TLS:     k.tls,
			SASLMechanism: plain.Mechanism{
				Username: k.Username,
				Password: k.Password,
			},
		}
	} else {
		dialer = &kafka.Dialer{
			Timeout:       k.Timeout(),
			TLS:           k.tls,
			SASLMechanism: nil,
		}
	}

	ctx, cancel := context.WithTimeout(k.Context, k.Timeout())
	defer cancel()

	conn, err := dialer.DialContext(ctx, "tcp", k.Host)
	if err != nil {
		return false, err.Error()
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(k.Timeout()))

	partitions, err := conn.ReadPartitions()
	if err != nil {
		return false, err.Error()
	}

	m := map[string]struct{}{}

	for _, p := range partitions {
		m[p.Topic] = struct{}{}
	}
	for t := range m {
		log.Debugf("[%s / %s / %s] Topic Name - %s", k.ProbeKind, k.ProbeName, k.ProbeTag, t)
	}

	return true, "Check Kafka Server Successfully!"

}
