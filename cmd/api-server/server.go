package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

var errHTTPServerFailed = errors.New("HTTP server failed")

func (app *Application) StartHTTPServer() {
	app.server = &http.Server{
		Addr:    app.config.BindAddress,
		Handler: app.routes(),
	}

	go func() {
		app.log.Infof("Starting HTTP server at %s", app.config.BindAddress)
		if err := app.server.ListenAndServe(); err != nil {
			app.eventChan <- fmt.Errorf("%w to start: %v", errHTTPServerFailed, err)
		} else {
			app.log.Info("Stopped HTTP server")
		}

	}()
}

func (app *Application) StopHTTPServer() {
	_ = app.server.Shutdown(context.Background())
}
