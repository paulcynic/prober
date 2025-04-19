package conf

import (
	"errors"
	"io"
	"net/http"
	httpClient "net/http"
	"os"
	"testing"
	"time"

	"github.com/megaease/easeprobe/monkey"
	"github.com/megaease/easeprobe/probe/client"
	clientConf "github.com/megaease/easeprobe/probe/client/conf"
	httpProbe "github.com/megaease/easeprobe/probe/http"
	"github.com/megaease/easeprobe/probe/tcp"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func testisExternalURL(url string, expects bool, t *testing.T) {
	if got := isExternalURL(url); got != expects {
		t.Errorf("isExternalURL(\"%s\") = %v, expected %v", url, got, expects)
	}
}

func TestPathAndURL(t *testing.T) {
	testisExternalURL("/tmp", false, t)
	testisExternalURL("//tmp", false, t)
	testisExternalURL("file:///tmp", false, t)
	testisExternalURL("http://", false, t)
	testisExternalURL("https://", false, t)
	testisExternalURL("hTtP://", false, t)
	testisExternalURL("http", false, t)
	testisExternalURL("https", false, t)
	testisExternalURL("ftp", false, t)
	testisExternalURL("hTtP://127.0.0.1", true, t)
	testisExternalURL("localhost", false, t)
	testisExternalURL("ftp://127.0.0.1", false, t)
}

func testScheduleYaml(t *testing.T, name string, sch Schedule, good bool) {
	var s Schedule
	err := yaml.Unmarshal([]byte(name), &s)
	if good {
		assert.Nil(t, err)
		assert.Equal(t, sch, s)
	} else {
		assert.NotNil(t, err)
	}

	buf, err := yaml.Marshal(sch)
	if good {
		assert.Nil(t, err)
		assert.Equal(t, name+"\n", string(buf))
	} else {
		assert.NotNil(t, err)
	}
}
func TestScheduleYaml(t *testing.T) {
	testScheduleYaml(t, "minutely", Minutely, true)
	testScheduleYaml(t, "hourly", Hourly, true)
	testScheduleYaml(t, "daily", Daily, true)
	testScheduleYaml(t, "weekly", Weekly, true)
	testScheduleYaml(t, "monthly", Monthly, true)
	testScheduleYaml(t, "none", None, true)
	testScheduleYaml(t, "yearly", 100, false)
	testScheduleYaml(t, "- bad", 100, false)
}

func TestGetYamlFileFromFile(t *testing.T) {
	if _, err := getYamlFileFromFile("/tmp/nonexistent"); err == nil {
		t.Errorf("getYamlFileFromFile(\"/tmp/nonexistent\") = nil, expected error")
	}

	tmpfile, err := os.CreateTemp("", "invalid*.yaml")
	if err != nil {
		t.Errorf("TempFile(\"invalid*.yaml\") %v", err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	// test empty file
	data, err := getYamlFileFromFile(tmpfile.Name())
	if err != nil {
		t.Errorf("getYamlFileFromFile(\"%s\") = %v, expected nil", tmpfile.Name(), err)
	}

	//confirm we read empty data
	if string(data) != "" {
		t.Errorf("getYamlFileFromFile(\"%s\") got data %s, expected nil", tmpfile.Name(), data)
	}
}

func TestGetYamlFileFromPath(t *testing.T) {
	tmpDir := t.TempDir()
	defer os.RemoveAll(tmpDir)

	// test empty dir
	_, err := getYamlFileFromFile(tmpDir)
	assert.NotNil(t, err)
}

const confVer = "version: 0.1.0\n"

const confHTTP = `
http:
  - name: dummy
    url: http://localhost:12345/dummy
    channels:
      - "telegram#Dev"
  - name: Local mTLS test
    url: https://localhost:8443/hello
    ca: ../mTLS/certs/ca.crt
    cert: ../mTLS/certs/client.b.crt
    key: ../mTLS/certs/client.b.key
  - name: MegaCloud
    url: https://cloud.megaease.cn/
    timeout: 2m
    interval: 30s
    channels:
      - "telegram#Dev"
  - name: Env Variables
    url: $WEB_SITE
`

func checkHTTPProbe(t *testing.T, probe httpProbe.HTTP) {
	switch probe.ProbeName {
	case "dummy":
		assert.Equal(t, probe.URL, "http://localhost:12345/dummy")
		assert.Equal(t, probe.Channels(), []string{"telegram#Dev"})
	case "Local mTLS test":
		assert.Equal(t, probe.URL, "https://localhost:8443/hello")
		assert.Equal(t, probe.CA, "../mTLS/certs/ca.crt")
		assert.Equal(t, probe.Cert, "../mTLS/certs/client.b.crt")
		assert.Equal(t, probe.Key, "../mTLS/certs/client.b.key")
	case "MegaCloud":
		assert.Equal(t, probe.URL, "https://cloud.megaease.cn/")
		assert.Equal(t, probe.ProbeTimeout, 2*time.Minute)
		assert.Equal(t, probe.ProbeTimeInterval, 30*time.Second)
		assert.Equal(t, probe.Channels(), []string{"telegram#Dev"})
	case "Env Variables":
		assert.Equal(t, probe.URL, os.Getenv("WEB_SITE"))
	default:
		t.Errorf("unexpected probe name %s", probe.ProbeName)
	}
}

const confTCP = `
tcp:
  - name: Example SSH
    host: example.com:22
    timeout: 10s
    interval: 2m
  - name: Example HTTP
    host: example.com:80
`

func checkTCPProbe(t *testing.T, probe tcp.TCP) {
	switch probe.ProbeName {
	case "Example SSH":
		assert.Equal(t, probe.Host, "example.com:22")
		assert.Equal(t, probe.ProbeTimeout, 10*time.Second)
		assert.Equal(t, probe.ProbeTimeInterval, 2*time.Minute)
	case "Example HTTP":
		assert.Equal(t, probe.Host, "example.com:80")
	default:
		t.Errorf("unexpected probe name %s", probe.ProbeName)
	}
}

const confClient = `
client:
  - name: Redis Native Client (local)
    driver: "redis"
    host: "localhost:6379"
    password: "abc123"
    channels:
      - test
  - name: MySQL Native Client (local)
    driver: "mysql"
    host: "localhost:3306"
    username: "root"
    password: "pass"
    ca: /home/chenhao/Github/mTLS/certs/ca.crt
    cert: /home/chenhao/Github/mTLS/certs/client.b.crt
    key: /home/chenhao/Github/mTLS/certs/client.b.key
`

func checkClientProbe(t *testing.T, probe client.Client) {
	switch probe.ProbeName {
	case "Redis Native Client (local)":
		assert.Equal(t, probe.DriverType, clientConf.Redis)
		assert.Equal(t, probe.Host, "localhost:6379")
		assert.Equal(t, probe.Password, "abc123")
		assert.Equal(t, probe.Channels(), []string{"test"})
	case "MySQL Native Client (local)":
		assert.Equal(t, probe.DriverType, clientConf.MySQL)
		assert.Equal(t, probe.Host, "localhost:3306")
		assert.Equal(t, probe.Username, "root")
		assert.Equal(t, probe.Password, "pass")
		assert.Equal(t, probe.CA, "/home/chenhao/Github/mTLS/certs/ca.crt")
		assert.Equal(t, probe.Cert, "/home/chenhao/Github/mTLS/certs/client.b.crt")
		assert.Equal(t, probe.Key, "/home/chenhao/Github/mTLS/certs/client.b.key")
	default:
		t.Errorf("unexpected probe name %s", probe.ProbeName)
	}
}

const confSettings = `
settings:
  name: "EaseProbeBot"
  icon: https://upload.wikimedia.org/wikipedia/commons/2/2d/Etcher-icon.png
  http:
    ip: 127.0.0.1
    port: 8181
    refresh: 5s
  probe:
    interval: 15s
  log:
    level: debug
    size: 1
  timeformat: "2006-01-02 15:04:05 UTC"
`

func checkSettings(t *testing.T, s Settings) {
	assert.Equal(t, s.Name, "EaseProbeBot")
	assert.Equal(t, s.IconURL, "https://upload.wikimedia.org/wikipedia/commons/2/2d/Etcher-icon.png")
	assert.Equal(t, s.HTTPServer.IP, "127.0.0.1")
	assert.Equal(t, s.HTTPServer.Port, "8181")
	assert.Equal(t, s.HTTPServer.AutoRefreshTime, 5*time.Second)
	assert.Equal(t, s.Probe.Interval, 15*time.Second)
	assert.Equal(t, s.Log.Level, LogLevel(log.DebugLevel))
	assert.Equal(t, s.Log.MaxSize, 1)
	assert.Equal(t, s.TimeFormat, "2006-01-02 15:04:05 UTC")
}

const confYAML = confVer + confHTTP + confClient + confSettings

func writeConfig(file, content string) error {
	return os.WriteFile(file, []byte(content), 0644)
}

func httpServer(port string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(confYAML))
	})
	mux.HandleFunc("/modified", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(confYAML + "  \n  \n"))
	})

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		panic(err)
	}
}

func TestConfig(t *testing.T) {
	file := "./config.yaml"
	err := writeConfig(file, confYAML)
	assert.Nil(t, err)

	// bad config
	os.Setenv("WEB_SITE", "\n - x::")
	_, err = New(&file)
	assert.NotNil(t, err)

	os.Setenv("WEB_SITE", "https://easeprobe.com")
	monkey.Patch(yaml.Marshal, func(v interface{}) ([]byte, error) {
		return nil, errors.New("marshal error")
	})
	_, err = New(&file)
	assert.Nil(t, err)
	monkey.UnpatchAll()

	_, err = New(&file)
	assert.Nil(t, err)
	conf := Get()

	assert.Equal(t, "EaseProbeBot", conf.Settings.Name)
	assert.Equal(t, "0.1.0", conf.Version)

	for _, v := range conf.HTTP {
		checkHTTPProbe(t, v)
	}
	for _, v := range conf.TCP {
		checkTCPProbe(t, v)
	}
	for _, v := range conf.Client {
		checkClientProbe(t, v)
	}
	checkSettings(t, conf.Settings)

	conf.InitAllLogs()
	probers := conf.AllProbers()
	assert.Equal(t, 9, len(probers))

	go httpServer("65535")
	url := "http://localhost:65535"
	os.Setenv("HTTP_AUTHORIZATION", "Basic dXNlcm5hbWU6cGFzc3dvcmQ=")
	os.Setenv("HTTP_TIMEOUT", "10")
	httpConf, err := New(&url)
	assert.Nil(t, err)
	assert.Equal(t, "EaseProbeBot", httpConf.Settings.Name)
	assert.Equal(t, "0.1.0", httpConf.Version)

	// test config modification
	assert.False(t, IsConfigModified(url))
	assert.False(t, IsConfigModified(url))
	url += "/modified"
	assert.True(t, IsConfigModified(url))

	probers = conf.AllProbers()
	assert.Equal(t, 9, len(probers))

	os.RemoveAll(file)
	os.RemoveAll("data")

	// error test
	url = "http://localhost:65534"
	_, err = New(&url)
	assert.NotNil(t, err)

	os.Setenv("HTTP_TIMEOUT", "invalid")
	_, err = New(&url)
	assert.NotNil(t, err)

	monkey.Patch(httpClient.NewRequest, func(method, url string, body io.Reader) (*http.Request, error) {
		return nil, errors.New("error")
	})
	url = "http://localhost"
	_, err = New(&url)
	assert.NotNil(t, err)

	monkey.UnpatchAll()
}

func TestEmptyProbes(t *testing.T) {
	myConf := confVer
	file := "./config.yaml"
	err := writeConfig(file, myConf)
	assert.Nil(t, err)

	conf, err := New(&file)
	assert.Nil(t, err)
	probers := conf.AllProbers()
	assert.Equal(t, 0, len(probers))

	os.RemoveAll(file)
	os.RemoveAll("data")
}

func TestFileConfigModificaiton(t *testing.T) {
	file := "./config.yaml"
	err := writeConfig(file, confYAML)
	assert.Nil(t, err)
	ResetPreviousYAMLFile()
	assert.False(t, IsConfigModified(file))
	assert.False(t, IsConfigModified(file))

	err = writeConfig(file, confYAML+"  \n\n")
	assert.Nil(t, err)
	assert.True(t, IsConfigModified(file))
	assert.False(t, IsConfigModified(file))

	err = writeConfig(file, confYAML+"\ninvalid")
	assert.Nil(t, err)
	assert.False(t, IsConfigModified(file))

	os.RemoveAll(file)
	assert.False(t, IsConfigModified(file))

}
