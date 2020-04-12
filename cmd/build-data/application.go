package main

import (
	"fmt"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"simplon.biz/corona/pkg/keystorage"
	"simplon.biz/corona/pkg/tokens"
)

type Application struct {
	log          *logrus.Logger
	tokenManager tokens.TokenManager
	keyStorage   keystorage.KeyStorage
	dumpPath     string
}

func NewApplication(log *logrus.Logger, tokenManager tokens.TokenManager, keyStorage keystorage.KeyStorage, dumpPath string) *Application {
	return &Application{
		log:          log,
		tokenManager: tokenManager,
		keyStorage:   keyStorage,
		dumpPath:     dumpPath,
	}
}

func (app *Application) DumpData() int {
	fullDump, err := NewRecordWriter(app.log, filepath.Join(app.dumpPath, "all.csv"))
	if err != nil {
		fullDump.log.Fatalf("Error creating dump file: %v", err)
	}
	fullDump.log.Info("Will dump full data here")
	defer fullDump.Abort()

	dayDumpers := make(map[string]*RecordWriter)
	defer func() {
		for _, d := range dayDumpers {
			d.Abort()
		}
	}()
	errors := make(chan interface{})
	records := make(chan keystorage.RawKeyRecord, 10)

	go app.keyStorage.ListRecords(records, errors)

	recordsProcessed := 0
out:
	for {
		select {
		case err := <-errors:
			app.log.Fatal(err)

		case record, more := <-records:
			if !more {
				break out
			}

			if err = fullDump.Write(&record); err != nil {
				fullDump.log.Fatalf("Error writing to file: %v", err)
			}

			dayDumper, ok := dayDumpers[record.DayNumber]
			if !ok {
				dayDumper, err = NewRecordWriter(app.log, filepath.Join(app.dumpPath, fmt.Sprintf("%s.csv", record.DayNumber)))
				if err != nil {
					dayDumper.log.Fatalf("Error creating dump file: %v", err)
				}
				dayDumper.log.Infof("Will dump data for day number %v here", record.DayNumber)
				dayDumpers[record.DayNumber] = dayDumper
			}

			if err = dayDumper.Write(&record); err != nil {
				dayDumper.log.Fatalf("Error writing to file: %v", err)
			}

			recordsProcessed += 1
		}
	}

	fullDump.log.Debug("Closing file")
	if err = fullDump.Close(); err != nil {
		fullDump.log.Error(err)
	}

	for _, d := range dayDumpers {
		d.log.Debug("Closing file")
		if err = d.Close(); err != nil {
			d.log.Error(err)
		}
	}

	return recordsProcessed
}
