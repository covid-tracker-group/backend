package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"simplon.biz/corona/pkg/config"
	"simplon.biz/corona/pkg/keystorage"
	"simplon.biz/corona/pkg/tokens"
)

var dataPath = flag.String("data", fmt.Sprintf("/var/lib/%s", config.AppName), "Directory to store all data")

func main() {
	flag.Parse()

	log := logrus.StandardLogger()
	log.SetLevel(logrus.DebugLevel)

	tokenManager, err := tokens.NewDiskTokenManager(filepath.Join(*dataPath, "tokens"))
	if err != nil {
		log.Fatalf("Can not create token manager: %v", err)
	}

	keyStorage, err := keystorage.NewDiskKeyStorage(filepath.Join(*dataPath, "records"))
	if err != nil {
		log.Fatalf("Can not create key storage: %v", err)
	}

	expireDailyTracingAuthorisationCodes(tokenManager, keyStorage)
}

func expireDailyTracingAuthorisationCodes(tokenManager tokens.TokenManager, keyStorage keystorage.KeyStorage) {
	expired := 0
	errors := 0
	tokenManager.Expire(func(token string) {
		contextLog := logrus.WithField("token", token)
		contextLog.Info("Expiring daily tracing authorisation code")
		if err := keyStorage.PurgeRecords(token); err != nil {
			contextLog.WithField("error", err).Error("Error remove keys for token")
			errors += 1
		} else {
			expired += 1
		}
	})

	logrus.WithFields(logrus.Fields{
		"expired": expired,
		"errors":  errors,
	}).Info("Finished expiring authorisation codes")
}
