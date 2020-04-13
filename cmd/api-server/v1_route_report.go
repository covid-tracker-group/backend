package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"simplon.biz/corona/pkg/authz"
)

type reportRequest struct {
	AuthorisationCode string      `json:"authorisation"`
	Keys              []keyRecord `json:"keys"`
}

type reportResponse struct {
	AuthorisationCode string `json:"authorisation"`
}

func (app *Application) report(w http.ResponseWriter, r *http.Request) {
	log := getRequestLog(r)

	var request reportRequest
	err := decodeJSONBody(w, r, &request)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			log.Errorf("Error decoding data: %v", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	log = log.WithField("medicalAuthCode", request.AuthorisationCode)
	_, err = app.authzManager.ValidateMedicalAuthCode(request.AuthorisationCode)
	if err != nil {
		var invErr authz.MedicalAuthCodeError
		if errors.As(err, &invErr) {
			log.WithField("error", invErr.Error()).Error("Invalid medical authorisation code")
			http.Error(w, invErr.Error(), http.StatusBadRequest)
			return
		}
		log.WithField("error", err).Error("Error validating medical authorisation code")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	tracingAuthCode, err := app.tokenManager.CreateToken()
	if err != nil {
		log.WithField("error", err).Error("Error creating tracing auth code")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if len(request.Keys) > 0 {
		if err = app.storeRecords(tracingAuthCode, request.Keys); err != nil {
			log.WithField("error", err).Error("Error adding key records")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	response := reportResponse{
		AuthorisationCode: tracingAuthCode,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.WithField("error", err).Error("Error encoding response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
