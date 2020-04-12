package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"simplon.biz/corona/pkg/keystorage"
	"simplon.biz/corona/pkg/tokens"
)

const appName = "covid-tracker"

var dataPath = flag.String("data", fmt.Sprintf("/var/lib/%s", appName), "Directory to store all data")
var dumpPath = flag.String("dump", fmt.Sprintf("/var/cache/%s", appName), "Directory to generated dumps")

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

	fi, err := os.Stat(*dumpPath)
	if err != nil {
		log.WithField("path", *dumpPath).Fatal("Invalid dump path: directory does not exist")
	}
	if !fi.Mode().IsDir() {
		log.WithField("path", *dumpPath).Fatal("Invalid dump path: not a directory")
	}

	app := NewApplication(log, tokenManager, keyStorage, *dumpPath)
	recordsProcessed := app.DumpData()
	log.Infof("Processed %d daily tracking key records", recordsProcessed)
}
