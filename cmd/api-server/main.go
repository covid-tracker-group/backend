package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bdlm/log"
	"github.com/sirupsen/logrus"
)

const appName = "covid-tracker"

var bindAddress = flag.String("bind", "127.0.0.1:8080", "Address for HTTP server")
var dataPath = flag.String("data", fmt.Sprintf("/var/lib/%s", appName), "Directory to store all data")

func main() {
	flag.Parse()

	app := NewApplication(Configuration{
		BindAddress: *bindAddress,
		DataPath:    *dataPath,
	})
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
