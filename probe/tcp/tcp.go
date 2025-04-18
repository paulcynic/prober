
// Package tcp is the tcp probe package
package tcp

import (
	"fmt"
	"net"

	"github.com/megaease/easeprobe/global"
	"github.com/megaease/easeprobe/probe/base"
	log "github.com/sirupsen/logrus"
)

// TCP implements a config for TCP
type TCP struct {
	base.DefaultProbe `yaml:",inline"`
	Host              string `yaml:"host" json:"host" jsonschema:"required,format=hostname,title=Host,description=The host to probe"`
	Proxy             string `yaml:"proxy" json:"proxy,omitempty" jsonschema:"format=hostname,title=Proxy,description=The proxy to use"`
	NoLinger          bool   `yaml:"nolinger" json:"nolinger" jsonschema:"format=nolinger,title=Disable SO_LINGER,description=Disable SO_LINGER TCP flag, default=false"`
}

// Config HTTP Config Object
func (t *TCP) Config(gConf global.ProbeSettings) error {
	kind := "tcp"
	tag := ""
	name := t.ProbeName
	t.DefaultProbe.Config(gConf, kind, tag, name, t.Host, t.DoProbe)

	log.Debugf("[%s / %s] configuration: %+v", t.ProbeKind, t.ProbeName, *t)
	return nil
}

// DoProbe return the checking result
func (t *TCP) DoProbe() (bool, string) {
	conn, err := t.GetProxyConnection(t.Proxy, t.Host)
	status := true
	message := ""
	if err != nil {
		message = fmt.Sprintf("Error: %v", err)
		log.Errorf("[%s / %s] error: %v", t.ProbeKind, t.ProbeName, err)
		status = false
	} else {
		message = "TCP Connection Established Successfully!"
		if tcpCon, ok := conn.(*net.TCPConn); ok && !t.NoLinger {
			tcpCon.SetLinger(0)
		}
		defer conn.Close()
	}
	return status, message
}
