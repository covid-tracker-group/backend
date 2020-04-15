package main

import (
	"net/http"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"simplon.biz/corona/pkg/config"
	"simplon.biz/corona/pkg/tokens"
)

type Application struct {
	config                  Configuration
	testingAuthTokenManager tokens.TokenManager
	eventChan               chan interface{}
	log                     *logrus.Logger
	server                  *http.Server
}

func NewApplication(cfg Configuration) *Application {
	log := logrus.StandardLogger()

	testingAuthTokenManager, err := tokens.NewDiskTokenManager(filepath.Join(cfg.DataPath, "test-tokens"), config.ExpireDailyTracingTokensAfter)
	if err != nil {
		log.Fatalf("Can not create testig auth token manager: %v", err)
	}

	return &Application{
		config:                  cfg,
		testingAuthTokenManager: testingAuthTokenManager,
		eventChan:               make(chan interface{}, 16),
		log:                     log,
	}
}
