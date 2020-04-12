package main

import (
	"encoding/csv"

	"github.com/sirupsen/logrus"
	"simplon.biz/corona/pkg/keystorage"
	"simplon.biz/corona/pkg/tools"
)

type RecordWriter struct {
	file      *tools.SecureFile
	csvWriter *csv.Writer
	log       *logrus.Entry
}

func NewRecordWriter(log *logrus.Logger, path string) (*RecordWriter, error) {
	contextLog := log.WithField("path", path)
	file, err := tools.OpenSecureFile(path)
	if err != nil {
		return &RecordWriter{log: contextLog}, err
	}

	csvWriter := csv.NewWriter(file)
	if err = csvWriter.Write([]string{"day_number", "daily_tracing_key"}); err != nil {
		return &RecordWriter{log: contextLog}, err
	}

	writer := &RecordWriter{
		file:      file,
		csvWriter: csvWriter,
		log:       contextLog,
	}
	return writer, nil
}

func (wr *RecordWriter) Write(record *keystorage.RawKeyRecord) error {
	err := wr.csvWriter.Write([]string{record.DayNumber, record.DailyTracingKey})
	if err != nil {
		wr.file.Abort()
	}
	return err
}

func (wr *RecordWriter) Abort() {
	wr.file.Abort()
}

func (wr *RecordWriter) Close() error {
	wr.csvWriter.Flush()
	if err := wr.csvWriter.Error(); err != nil {
		wr.file.Abort()
		return err
	}
	return wr.file.Close()
}
