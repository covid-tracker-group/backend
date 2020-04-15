package main

import (
	"net/http"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"simplon.biz/corona/pkg/config"
	"simplon.biz/corona/pkg/keystorage"
	"simplon.biz/corona/pkg/tokens"
)

type Application struct {
	config                  Configuration
	eventChan               chan interface{}
	log                     *logrus.Logger
	server                  *http.Server
	testingAuthTokenManager tokens.TokenManager
	tracingAuthTokenManager tokens.TokenManager

	keyStorage keystorage.KeyStorage
}

func NewApplication(cfg Configuration) *Application {
	log := logrus.StandardLogger()

	testingAuthTokenManager, err := tokens.NewDiskTokenManager(filepath.Join(cfg.DataPath, "test-tokens"), config.ExpireDailyTracingTokensAfter)
	if err != nil {
		log.Fatalf("Can not create testig auth token manager: %v", err)
	}

	tracingAuthTokenManager, err := tokens.NewDiskTokenManager(filepath.Join(cfg.DataPath, "trace-tokens"), config.ExpireDailyTracingTokensAfter)
	if err != nil {
		log.Fatalf("Can not create tracing auth token manager: %v", err)
	}

	keyStorage, err := keystorage.NewDiskKeyStorage(filepath.Join(cfg.DataPath, "records"))
	if err != nil {
		log.Fatalf("Can not create key storage: %v", err)
	}

	return &Application{
		config:                  cfg,
		eventChan:               make(chan interface{}, 16),
		log:                     log,
		testingAuthTokenManager: testingAuthTokenManager,
		tracingAuthTokenManager: tracingAuthTokenManager,
		keyStorage:              keyStorage,
	}
}
