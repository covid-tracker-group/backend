package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"simplon.biz/corona/pkg/tokens"
	"simplon.biz/corona/pkg/tools"
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
	err := tools.DecodeJSONBody(w, r, &request)
	if err != nil {
		var mr *tools.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Message, mr.Status)
		} else {
			log.WithError(err).Error("Error decoding data")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	log = log.WithField("medicalAuthCode", request.AuthorisationCode)
	ok, err := app.testingAuthTokenManager.VerifyToken(request.AuthorisationCode)
	if err != nil {
		log.WithError(err).Error("Error verifying health test authorisation code")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if !ok {
		log.Error("Invalid health test authorisation code")
		http.Error(w, "Invalid authorisation code", http.StatusBadRequest)
		return
	}
	if err = app.testingAuthTokenManager.RetractToken(request.AuthorisationCode); err != nil {
		log.WithError(err).Error("Error retracting health test authorisation code")
	}

	tracingAuthCode := tokens.NewTracingAuthenticationToken()
	if err = app.tracingAuthTokenManager.StoreToken(tracingAuthCode); err != nil {
		log.WithError(err).Error("Error creating tracing auth code")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if len(request.Keys) > 0 {
		if err = app.storeRecords(tracingAuthCode.GetCode(), request.Keys); err != nil {
			log.WithError(err).Error("Error adding key records")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	response := reportResponse{
		AuthorisationCode: tracingAuthCode.GetCode(),
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.WithError(err).Error("Error encoding response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
