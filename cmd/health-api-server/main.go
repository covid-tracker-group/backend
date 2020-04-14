package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/bdlm/log"
	"github.com/sirupsen/logrus"
	"simplon.biz/corona/pkg/config"
)

var bindAddress = flag.String("bind", "127.0.0.1:8080", "Address for HTTP server")
var proxyURL = flag.String("proxy", "", "URL to proxy web app requests to")
var httpPath = flag.String("httpData", fmt.Sprintf("/var/lib/%s/http", config.AppName), "Folder with web app")

func main() {
	flag.Parse()

	var err error
	config := Configuration{
		BindAddress: *bindAddress,
		HttpPath:    *httpPath,
	}
	if *proxyURL != "" {
		config.ProxyURL, err = url.Parse(*proxyURL)
		if err != nil {
			log.WithField("proxy", *proxyURL).Fatalf("Invalid proxy URL: %v", err)
		}
	} else {
		st, err := os.Stat(*httpPath)
		if err != nil || !st.Mode().IsDir() {
			log.WithField("httpPath", *httpPath).Fatalf("Invalid HTTP path: %v", err)
		}
	}
	app := NewApplication(config)
	app.log.SetLevel(logrus.DebugLevel)

	abortSignal := make(chan os.Signal, 4)

	signal.Notify(abortSignal, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)

	app.StartHTTPServer()

out:
	for {
		select {
		case event := <-app.eventChan:
			switch evt := event.(type) {
			case error:
				log.Fatal(evt)
			}

		case sgn := <-abortSignal:
			log.Infof("Stopping all activity on %s", sgn.String())
			app.StopHTTPServer()
			break out
		}
	}
}
