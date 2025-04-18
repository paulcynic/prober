// Package web is the web server of easeprobe.
package web

import (
	"context"
	"errors"
	"html"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/megaease/easeprobe/conf"
	"github.com/megaease/easeprobe/global"
	"github.com/megaease/easeprobe/probe"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
)

var probers *[]probe.Prober
var webServer *http.Server

func getRefreshInterval(refersh string) time.Duration {
	interval := conf.Get().Settings.HTTPServer.AutoRefreshTime
	if strings.TrimSpace(refersh) == "" {
		return interval
	}
	r, err := time.ParseDuration(refersh)
	if err != nil {
		log.Errorf("[Web] Invalid refresh time: %s", err)
		return interval
	}
	return r
}

func getStatus(status string) *probe.Status {
	if status == "" {
		return nil
	}
	var s probe.Status
	s.Status(status)
	return &s
}

func toFloat(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}
func toInt(str string) (int, error) {
	return strconv.Atoi(str)
}
func getNum[T any](str string, _default T, convert func(string) (T, error)) T {
	if str == "" {
		return _default
	}
	n, err := convert(str)
	if err != nil {
		log.Debugf("[Web] Invalid number value: %s", err)
		return _default
	}
	return n
}
func getStr(str string) string {
	return strings.TrimSpace(html.EscapeString(str))
}


// SetProbers set the probers
func SetProbers(p []probe.Prober) {
	probers = &p
}

// Server is the http server
func Server() {

	c := conf.Get()
	host := c.Settings.HTTPServer.IP
	port := c.Settings.HTTPServer.Port

	// Configure the http server
	if len(host) > 0 && net.ParseIP(host) == nil {
		host = global.DefaultHTTPServerIP
	}
	p, err := strconv.Atoi(port)
	if err != nil || p <= 1024 || p > 65535 {
		log.Warnf("[Web] Invalid port number: %s, use the default value: %s", port, global.DefaultHTTPServerPort)
		port = global.DefaultHTTPServerPort
	} else {
		port = c.Settings.HTTPServer.Port
	}

	// Configure the auto refresh time of the SLA page
	if c.Settings.HTTPServer.AutoRefreshTime == 0 {
		c.Settings.HTTPServer.AutoRefreshTime = global.DefaultProbeInterval
	}
	log.Debugf("[Web] Auto refresh interval time: %s", c.Settings.HTTPServer.AutoRefreshTime)

	// Prepare the router
	r := chi.NewRouter()

	filename := c.Settings.HTTPServer.AccessLog.File
	if len(filename) > 0 {
		log.Infof("[Web] Access Log output file: %s", filename)
		logger := c.Settings.HTTPServer.AccessLog.Logger
		r.Use(NewStructuredLogger(logger))
	}

	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.StripSlashes)

	r.Get("/metrics", promhttp.Handler().ServeHTTP)

	server, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		log.Fatalf("[Web] Failed to start the http server: %s", err)
	}
	log.Infof("[Web] HTTP server is listening on %s:%s", host, port)

	// Start the http server
	go func() {
		webServer = &http.Server{Handler: r}
		if err := webServer.Serve(server); !errors.Is(err, http.ErrServerClosed) {
			log.Errorf("[Web] HTTP server error: %s", err)
		}
		log.Info("[Web] HTTP server is stopped.")
	}()

}

// Shutdown the http server
func Shutdown() {
	if webServer == nil {
		log.Debugf("[Web] HTTP server is not running, skip to shutdown")
		return
	}
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), global.DefaultTimeOut)
	defer shutdownRelease()

	if err := webServer.Shutdown(shutdownCtx); err != nil {
		log.Errorf("[Web] Failed to shutdown the http server: %s", err)
	}
	log.Info("[Web] HTTP server is shutdown")
	webServer = nil

}
