package main

import (
	"github.com/sirupsen/logrus"
	"simplon.biz/corona/pkg/keystorage"
	"simplon.biz/corona/pkg/tools"
)

type RecordWriter struct {
	log *logrus.Entry
	*tools.SecureCSVWriter
}

func NewRecordWriter(log *logrus.Logger, path string) (writer *RecordWriter, err error) {
	writer = &RecordWriter{
		log: log.WithField("path", path),
	}
	writer.SecureCSVWriter, err = tools.NewSecureCSVWriter(path)
	if err == nil {
		err = writer.SecureCSVWriter.Write([]string{"day_number", "daily_tracing_key"})
		if err != nil {
			writer.SecureCSVWriter.Abort()
		}
	}
	return
}

func (rw *RecordWriter) Write(record *keystorage.RawKeyRecord) error {
	return rw.SecureCSVWriter.Write([]string{record.DayNumber, record.DailyTracingKey})
}

func (rw *RecordWriter) Abort() {
	rw.SecureCSVWriter.Abort()
}

func (rw *RecordWriter) Close() error {
	return rw.SecureCSVWriter.Close()
}
