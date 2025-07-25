// Package main is the entry point for the easeprobe command.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/megaease/easeprobe/conf"
	"github.com/megaease/easeprobe/daemon"
	"github.com/megaease/easeprobe/global"
	"github.com/megaease/easeprobe/metric"
	"github.com/megaease/easeprobe/probe"
	"github.com/megaease/easeprobe/web"
	"github.com/prometheus/client_golang/prometheus"

	log "github.com/sirupsen/logrus"
)

func getEnvOrDefault(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func showVersion() {

	var v = global.Ver

	bi, ok := debug.ReadBuildInfo()
	if !ok {
		v = fmt.Sprintf("%v %v", global.DefaultProg, v)
		fmt.Println(v)
		return
	}

	for _, s := range bi.Settings {
		switch s.Key {
		case "vcs.revision":
			v = fmt.Sprintf("%v %v", v, s.Value[:9])
		case "vcs.time":
			v = fmt.Sprintf("%v %v", v, s.Value)
		}
	}

	v = fmt.Sprintf("%v %v %v", global.DefaultProg, v, bi.GoVersion)
	fmt.Println(v)
}

func main() {
	////////////////////////////////////////////////////////////////////////////
	//          Parse command line arguments and config file settings         //
	////////////////////////////////////////////////////////////////////////////

	dryNotify := flag.Bool("d", os.Getenv("PROBE_DRY") == "true", "dry notification mode")
	yamlFile := flag.String("f", getEnvOrDefault("PROBE_CONFIG", "config.yaml"), "configuration file")
	version := flag.Bool("v", false, "prints version")
	flag.Parse()

	if *version {
		showVersion()
		os.Exit(0)
	}

	// Create metrics for config file
	metrics := conf.NewConfigMetrics()
	c, err := conf.New(yamlFile)
	if err != nil {
		metrics.Status.With(metric.AddConstLabels(prometheus.Labels{
			"endpoint": *yamlFile,
		}, prometheus.Labels{})).Set(float64(0))
		log.Errorln("Fatal: Cannot read the YAML configuration file!")
		os.Exit(-1)
	}
	metrics.Status.With(metric.AddConstLabels(prometheus.Labels{
		"endpoint": *yamlFile,
	}, prometheus.Labels{})).Set(float64(1))
	metrics.Timestamp.With(metric.AddConstLabels(prometheus.Labels{
		"endpoint": *yamlFile,
	}, prometheus.Labels{})).Set(float64(time.Now().Unix()))
	currentTS := time.Now().Unix()
	log.Infof("Current timestamp: %d (%s)", currentTS, time.Unix(currentTS, 0).Format(time.RFC3339))

	// Create the pid file if the file name is not empty
	c.Settings.PIDFile = strings.TrimSpace(c.Settings.PIDFile)
	if len(c.Settings.PIDFile) > 0 && c.Settings.PIDFile != "-" {
		d, err := daemon.NewPIDFile(c.Settings.PIDFile)
		if err != nil {
			log.Errorf("Fatal: Cannot create the PID file: %s!", err)
			os.Exit(-1)
		}
		log.Infof("Successfully created the PID file: %s", d.PIDFile)
		defer d.RemovePIDFile()
	} else {
		if len(c.Settings.PIDFile) == 0 {
			log.Info("Skipping PID file creation (pid file is empty).")
		} else {
			log.Info("Skipping PID file creation (pid file is set to '-').")
		}
	}

	c.InitAllLogs()

	// if dry notification mode is specified in command line, overwrite the configuration
	if *dryNotify {
		log.Infoln("Dry Notification Mode...")
	}
	////////////////////////////////////////////////////////////////////////////
	//                          Start the HTTP Server                         //
	////////////////////////////////////////////////////////////////////////////
	// if error happens, the EaseProbe will exit
	web.Server()

	////////////////////////////////////////////////////////////////////////////
	//                  Configure all of Probers and Notifiers                //
	////////////////////////////////////////////////////////////////////////////

	// Probers
	probers := c.AllProbers()
	// Configure the Probes
	probers = configProbers(probers)
	if len(probers) == 0 {
		log.Fatal("No probes configured, exiting...")
	}

	////////////////////////////////////////////////////////////////////////////
	//                          Start the EaseProbe                           //
	////////////////////////////////////////////////////////////////////////////

	// wait group for probers
	var wg sync.WaitGroup
	// the exit channel for all probers
	doneProbe := make(chan bool, len(probers))
	// 2) Start the Probers
	doneSave := make(chan bool)
	// the channel for saving the probe result data
	saveChannel := make(chan probe.Result, len(probers))

	runProbers(probers, &wg, doneProbe, saveChannel)
	// 4) Set probers into web server
	web.SetProbers(probers)

	////////////////////////////////////////////////////////////////////////////
	//                          Rotate the log file                           //
	////////////////////////////////////////////////////////////////////////////
	rotateLog := make(chan os.Signal, 1)
	doneRotate := make(chan bool, 1)
	signal.Notify(rotateLog, syscall.SIGHUP)
	go func() {
		for {
			c := conf.Get()
			select {
			case <-doneRotate:
				log.Info("Received the exit signal, Rotating log file process exiting...")
				c.Settings.Log.Close()
				c.Settings.HTTPServer.AccessLog.Close()
				return
			case <-rotateLog:
				log.Info("Received SIGHUP, rotating the log file...")
				c.Settings.Log.Rotate()
				c.Settings.HTTPServer.AccessLog.Rotate()
			}
		}
	}()

	////////////////////////////////////////////////////////////////////////////
	//                         Graceful Shutdown / Re-Run                     //
	////////////////////////////////////////////////////////////////////////////
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGTERM)

	// the graceful shutdown process
	exit := func() {
		web.Shutdown()
		for i := 0; i < len(probers); i++ {
			if probers[i].Result().Status != probe.StatusBad {
				doneProbe <- true
			}
		}
		wg.Wait()
		doneSave <- true
		doneRotate <- true
	}

	// the graceful restart process
	reRun := func() {
		exit()
		p, e := os.StartProcess(os.Args[0], os.Args, &os.ProcAttr{
			Env:   os.Environ(),
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		})
		if e != nil {
			log.Errorf("!!! FAILED TO RESTART THE EASEPROBE: %v !!!", e)
			return
		}
		log.Infof("!!! RESTART THE EASEPROBE SUCCESSFULLY - PID=[%d] !!!", p.Pid)
	}

	// Monitor the configuration file
	monConf := make(chan bool, 1)
	go monitorYAMLFile(*yamlFile, monConf, metrics)

	// wait for the exit and restart signal
	select {
	case <-done:
		log.Info("!!! RECEIVED THE SIGTERM EXIT SIGNAL, EXITING... !!!")
		exit()
	case <-monConf:
		log.Info("!!! RECEIVED THE RESTART EVENT, RESTARTING... !!!")
		reRun()
	}

	log.Info("Graceful Exit Successfully!")
}

func monitorYAMLFile(path string, monConf chan bool, metrics *conf.ConfigMetrics) {
	for {
		if conf.IsConfigModified(path, metrics) {
			log.Infof("The configuration file [%s] has been modified, restarting...", path)
			// Set metrics after YAML modification
			metrics.Timestamp.With(metric.AddConstLabels(prometheus.Labels{
				"endpoint": path,
			}, prometheus.Labels{})).Set(float64(time.Now().Unix()))
			currentTS := time.Now().Unix()
			log.Infof("Current timestamp: %d (%s)", currentTS, time.Unix(currentTS, 0).Format(time.RFC3339))
			monConf <- true
			break
		}
		log.Debugf("The configuration file [%s] has not been modified", path)
		time.Sleep(global.DefaultConfigFileCheckInterval)
	}
}
