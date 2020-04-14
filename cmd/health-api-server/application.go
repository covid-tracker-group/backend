package main

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type Application struct {
	config    Configuration
	eventChan chan interface{}
	log       *logrus.Logger
	server    *http.Server
}

func NewApplication(config Configuration) *Application {
	log := logrus.StandardLogger()

	return &Application{
		config:    config,
		eventChan: make(chan interface{}, 16),
		log:       log,
	}
}
