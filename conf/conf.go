// Package conf is the configuration of the application
package conf

import (
	"bytes"
	"io"
	httpClient "net/http"
	netUrl "net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/megaease/easeprobe/global"
	"github.com/megaease/easeprobe/probe"
	"github.com/megaease/easeprobe/probe/client"
	"github.com/megaease/easeprobe/probe/http"
	"github.com/megaease/easeprobe/probe/tcp"
	"github.com/megaease/easeprobe/probe/tls"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var config *Conf

// Get return the global configuration
func Get() *Conf {
	return config
}

// Schedule is the schedule.
type Schedule int

// Schedule enum
const (
	None Schedule = iota
	Minutely
	Hourly
	Daily
	Weekly
	Monthly
)

var scheduleToString = map[Schedule]string{
	Minutely: "minutely",
	Hourly:   "hourly",
	Daily:    "daily",
	Weekly:   "weekly",
	Monthly:  "monthly",
	None:     "none",
}

var stringToSchedule = global.ReverseMap(scheduleToString)

// MarshalYAML marshal the configuration to yaml
func (s Schedule) MarshalYAML() (interface{}, error) {
	return global.EnumMarshalYaml(scheduleToString, s, "Schedule")
}

// UnmarshalYAML is unmarshal the debug level
func (s *Schedule) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return global.EnumUnmarshalYaml(unmarshal, stringToSchedule, s, None, "Schedule")
}

// Probe is the settings of prober
type Probe struct {
	Interval                             time.Duration `yaml:"interval" json:"interval,omitempty" jsonschema:"type=string,format=duration,title=Probe Interval,description=the interval of probe,default=1m"`
	Timeout                              time.Duration `yaml:"timeout" json:"timeout,omitempty" jsonschema:"type=string,format=duration,title=Probe Timeout,description=the timeout of probe,default=30s"`
	global.StatusChangeThresholdSettings `yaml:",inline" json:",inline"`
	global.NotificationStrategySettings  `yaml:"alert" json:"alert" jsonschema:"title=Alert,description=the alert settings"`
}

// HTTPServer is the settings of http server
type HTTPServer struct {
	IP              string        `yaml:"ip" json:"ip" jsonschema:"title=Web Server IP,description=the local ip address of the http server need to listen on,example=0.0.0.0"`
	Port            string        `yaml:"port" json:"port" jsonschema:"type=integer,title=Web Server Port,description=port of the http server,default=8181"`
	AutoRefreshTime time.Duration `yaml:"refresh" json:"refresh,omitempty" jsonschema:"type=string,title=Auto Refresh Time,description=auto refresh time of the http server,example=5s"`
	AccessLog       Log           `yaml:"log" json:"log,omitempty" jsonschema:"title=Access Log,description=access log of the http server"`
}

// Settings is the EaseProbe configuration
type Settings struct {
	Name       string     `yaml:"name" json:"name,omitempty" jsonschema:"title=EaseProbe Name,description=The name of the EaseProbe instance,default=EaseProbe"`
	IconURL    string     `yaml:"icon" json:"icon,omitempty" jsonschema:"title=Icon URL,description=The URL of the icon of the EaseProbe instance"`
	PIDFile    string     `yaml:"pid" json:"pid,omitempty" jsonschema:"title=PID File,description=The PID file of the EaseProbe instance ('' or '-' means no PID file)"`
	Log        Log        `yaml:"log" json:"log,omitempty" jsonschema:"title=EaseProbe Log,description=The log settings of the EaseProbe instance"`
	TimeFormat string     `yaml:"timeformat" json:"timeformat,omitempty" jsonschema:"title=Time Format,description=The time format of the EaseProbe instance,default=2006-01-02 15:04:05Z07:00"`
	TimeZone   string     `yaml:"timezone" json:"timezone,omitempty" jsonschema:"title=Time Zone,description=The time zone of the EaseProbe instance,example=Asia/Shanghai,example=Europe/Berlin,default=UTC"`
	Probe      Probe      `yaml:"probe" json:"probe,omitempty" jsonschema:"title=Probe Settings,description=The global probe settings of the EaseProbe instance"`
	HTTPServer HTTPServer `yaml:"http" json:"http,omitempty" jsonschema:"title=HTTP Server Settings,description=The HTTP server settings of the EaseProbe instance"`
}

// Conf is Probe configuration
type Conf struct {
	Version  string          `yaml:"version" json:"version,omitempty" jsonschema:"title=Version,description=Version of the EaseProbe configuration"`
	HTTP     []http.HTTP     `yaml:"http" json:"http,omitempty" jsonschema:"title=HTTP Probe,description=HTTP Probe Configuration"`
	TCP      []tcp.TCP       `yaml:"tcp" json:"tcp,omitempty" jsonschema:"title=TCP Probe,description=TCP Probe Configuration"`
	Client   []client.Client `yaml:"client" json:"client,omitempty" jsonschema:"title=Native Client Probe,description=Native Client Probe Configuration"`
	TLS      []tls.TLS       `yaml:"tls" json:"tls,omitempty" jsonschema:"title=TLS Probe,description=TLS Probe Configuration"`
	Settings Settings        `yaml:"settings" json:"settings,omitempty" jsonschema:"title=Global Settings,description=EaseProbe Global configuration"`
}


// Check if string is a url
func isExternalURL(url string) bool {
	if _, err := netUrl.ParseRequestURI(url); err != nil {
		log.Debugf("ParseRequestedURI: %s failed to parse with error %v", url, err)
		return false
	}

	parts, err := netUrl.Parse(url)
	if err != nil || parts.Host == "" || !strings.HasPrefix(parts.Scheme, "http") {
		log.Debugf("Parse: %s failed Scheme: %s, Host: %s (err: %v)", url, parts.Scheme, parts.Host, err)
		return false
	}

	return true
}

func getYamlFileFromHTTP(url string) ([]byte, error) {
	r, err := httpClient.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if os.Getenv("HTTP_AUTHORIZATION") != "" {
		r.Header.Set("Authorization", os.Getenv("HTTP_AUTHORIZATION"))
	}

	httpClientObject := httpClient.Client{}
	if os.Getenv("HTTP_TIMEOUT") != "" {
		timeout, err := strconv.ParseInt(os.Getenv("HTTP_TIMEOUT"), 10, 64)
		if err != nil {
			return nil, err
		}
		httpClientObject.Timeout = time.Duration(timeout) * time.Second
	}

	resp, err := httpClientObject.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func getYamlFileFromFile(path string) ([]byte, error) {
	f, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, err
	}
	if f.IsDir() {
		return mergeYamlFiles(path)
	}
	return os.ReadFile(path)
}

func getYamlFile(path string) ([]byte, error) {
	if isExternalURL(path) {
		return getYamlFileFromHTTP(path)
	}
	return getYamlFileFromFile(path)
}

// previousYAMLFile is the content of the configuration file
var previousYAMLFile []byte

// ResetPreviousYAMLFile resets the previousYAMLFile
func ResetPreviousYAMLFile() {
	previousYAMLFile = nil
}

// IsConfigModified checks if the configuration file is modified
func IsConfigModified(path string) bool {

	var content []byte
	var err error
	if isExternalURL(path) {
		content, err = getYamlFileFromHTTP(path)
	} else {
		content, err = getYamlFileFromFile(path)
	}

	if err != nil {
		log.Warnf("Failed to get the configuration file [%s]: %v", path, err)
		return false
	}

	// if it is the fisrt time to read the configuration file, we will not restart the program
	if previousYAMLFile == nil {
		previousYAMLFile = content
		return false
	}

	//  if the configuration file is invalid, we will not restart the program
	testConf := Conf{}
	err = yaml.Unmarshal(content, &testConf)
	if err != nil {
		log.Warnf("Invalid configuration file [%s]: %v", path, err)
		return false
	}

	// check if the configuration file is modified
	modified := !bytes.Equal(content, previousYAMLFile)
	previousYAMLFile = content
	return modified
}

// New read the configuration from yaml
func New(conf *string) (*Conf, error) {
	c := Conf{
		HTTP:   []http.HTTP{},
		TCP:    []tcp.TCP{},
		Client: []client.Client{},
		TLS:    []tls.TLS{},
		Settings: Settings{
			Name:       global.DefaultProg,
			IconURL:    global.DefaultIconURL,
			PIDFile:    filepath.Join(global.GetWorkDir(), global.DefaultPIDFile),
			Log:        NewLog(),
			TimeFormat: global.DefaultTimeFormat,
			TimeZone:   global.DefaultTimeZone,
			Probe: Probe{
				Interval: global.DefaultProbeInterval,
				Timeout:  global.DefaultTimeOut,
			},
			HTTPServer: HTTPServer{
				IP:        global.DefaultHTTPServerIP,
				Port:      global.DefaultHTTPServerPort,
				AccessLog: NewLog(),
			},
		},
	}
	y, err := getYamlFile(*conf)
	if err != nil {
		log.Errorf("error: %v ", err)
		return &c, err
	}

	y = []byte(os.ExpandEnv(string(y)))

	err = yaml.Unmarshal(y, &c)
	if err != nil {
		log.Errorf("error: %v", err)
		return &c, err
	}

	// Initialization
	c.Settings.Log.InitLog(nil)
	global.InitEaseProbeWithTime(c.Settings.Name, c.Settings.IconURL,
		c.Settings.TimeFormat, c.Settings.TimeZone)

	config = &c

	log.Infoln("Load the configuration file successfully!")
	if log.GetLevel() >= log.DebugLevel {
		s, err := yaml.Marshal(c)
		if err != nil {
			log.Debugf("%v\n%+v", err, c)
		} else {
			log.Debugf("\n%s", string(s))
		}
	}

	return &c, err
}

// InitAllLogs initialize all logs
func (conf *Conf) InitAllLogs() {

	conf.Settings.Log.InitLog(nil)
	conf.Settings.Log.LogInfo("Application")

	conf.Settings.HTTPServer.AccessLog.InitLog(log.New())
	conf.Settings.HTTPServer.AccessLog.LogInfo("Web Access")
}

// isProbe checks whether a interface is a probe type
func isProbe(t reflect.Type) bool {
	modelType := reflect.TypeOf((*probe.Prober)(nil)).Elem()
	return t.Implements(modelType)
}

// AllProbers return all probers
func (conf *Conf) AllProbers() []probe.Prober {
	log.Debugf("--------- Process the probers settings ---------")
	return allProbersHelper(*conf)
}

func allProbersHelper(i interface{}) []probe.Prober {

	var probers []probe.Prober
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	if t.Kind() != reflect.Struct {
		return probers
	}

	for i := 0; i < t.NumField(); i++ {
		tField := t.Field(i).Type.Kind()
		if tField == reflect.Struct {
			probers = append(probers, allProbersHelper(v.Field(i).Interface())...)
			continue
		}
		if tField != reflect.Slice {
			continue
		}

		vField := v.Field(i)
		for j := 0; j < vField.Len(); j++ {
			if !isProbe(vField.Index(j).Addr().Type()) {
				continue
			}

			log.Debugf("--> %s / %s / %+v", t.Field(i).Name, t.Field(i).Type.Kind(), vField.Index(j))
			probers = append(probers, vField.Index(j).Addr().Interface().(probe.Prober))
		}
	}

	return probers
}
