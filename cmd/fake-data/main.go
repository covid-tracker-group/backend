package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"simplon.biz/corona/pkg/keystorage"
	"simplon.biz/corona/pkg/tokens"
	"simplon.biz/corona/pkg/tools"
)

const appName = "covid-tracker"

var dataPath = flag.String("data", fmt.Sprintf("/var/lib/%s", appName), "Directory to store all data")

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

	tokens := make([]string, 0, 10000)
	for i := 0; i < 10000; i++ {
		token, err := tokenManager.CreateToken()
		if err != nil {
			log.Fatalf("Error creating token: %v", err)
		}
		tokens = append(tokens, token)
	}
	log.Infof("Generated %d authorisation tokens", len(tokens))

	now := time.Now()
	thePast := now.Add(time.Hour * -(24 * 15))
	var expiredTokens []string
	for i := 0; i < len(tokens)/10; i++ {
		expiredTokens = append(expiredTokens, tokens[i])
		path := filepath.Join(*dataPath, "tokens", tokens[i])
		_ = os.Chtimes(path, now, thePast)
	}
	log.Infof("Expired %d authorisation tokens", len(expiredTokens))

	var records []keystorage.KeyRecord
	toDayNumber := tools.TimeToDayNumber(now)

	recordsGenerated := 0
	tracingKey := make([]byte, 16)
	for _, token := range tokens {
		days := rand.Intn(7) + 1
		records = records[:0]
		for j := 0; j < days; j++ {
			_, _ = rand.Read(tracingKey)

			records = append(records, keystorage.KeyRecord{
				ProcessedAt:     now.Add(time.Hour * -time.Duration(rand.Intn(40))),
				DayNumber:       toDayNumber - rand.Intn(10),
				DailyTracingKey: base64.StdEncoding.EncodeToString(tracingKey),
			})
			recordsGenerated += 1

		}
		err = keyStorage.AddKeyRecords(token, records)
		if err != nil {
			log.Fatalf("Error writing key record: %v", err)
		}
	}
	log.Infof("Generated %d daily key records for %d authorisation tokens", recordsGenerated, len(tokens))
}
