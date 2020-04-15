package main

import (
	"net/http"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"simplon.biz/corona/pkg/authz"
	"simplon.biz/corona/pkg/config"
	"simplon.biz/corona/pkg/keystorage"
	"simplon.biz/corona/pkg/tokens"
)

type Application struct {
	config       Configuration
	eventChan    chan interface{}
	log          *logrus.Logger
	server       *http.Server
	tokenManager tokens.TokenManager
	keyStorage   keystorage.KeyStorage
	authzManager *authz.AuthorisationManager
}

func NewApplication(cfg Configuration) *Application {
	log := logrus.StandardLogger()

	tokenManager, err := tokens.NewDiskTokenManager(filepath.Join(cfg.DataPath, "tokens"), config.ExpireDailyTracingTokensAfter)
	if err != nil {
		log.Fatalf("Can not create token manager: %v", err)
	}

	keyStorage, err := keystorage.NewDiskKeyStorage(filepath.Join(cfg.DataPath, "records"))
	if err != nil {
		log.Fatalf("Can not create key storage: %v", err)
	}

	return &Application{
		config:       cfg,
		eventChan:    make(chan interface{}, 16),
		log:          log,
		tokenManager: tokenManager,
		keyStorage:   keyStorage,
		authzManager: authz.NewAuthorisationManager(),
	}
}
