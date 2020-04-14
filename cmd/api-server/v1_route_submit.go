package main

import (
	"errors"
	"net/http"
	"time"

	"simplon.biz/corona/pkg/keystorage"
	"simplon.biz/corona/pkg/tools"
)

type keyRecord struct {
	DayNumber       int    `json:"day_number"`
	DailyTracingKey string `json:"key"`
}

type submitRequest struct {
	Keys []keyRecord `json:"keys"`
}

func (app *Application) submit(w http.ResponseWriter, r *http.Request) {
	log := getRequestLog(r)

	var request submitRequest
	err := tools.DecodeJSONBody(w, r, &request)
	if err != nil {
		var mr *tools.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Message, mr.Status)
		} else {
			log.Errorf("Error decoding data: %v", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	tracingAuthCode := getTracingAuthenticationCode(r)
	if err = app.storeRecords(tracingAuthCode, request.Keys); err != nil {
		log.WithField("error", err).Error("Error adding key records")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *Application) storeRecords(tracingAuthCode string, keys []keyRecord) error {
	var records []keystorage.KeyRecord
	now := time.Now()
	for _, record := range keys {
		records = append(records, keystorage.KeyRecord{
			ProcessedAt:     now,
			DayNumber:       record.DayNumber,
			DailyTracingKey: record.DailyTracingKey,
		})
	}
	if len(records) > 0 {
		err := app.keyStorage.AddKeyRecords(tracingAuthCode, records)
		if err != nil {
			return err
		}
	}
	return nil
}
