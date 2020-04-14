package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"simplon.biz/corona/pkg/config"
	"simplon.biz/corona/pkg/tools"
)

type requestCodesRequest struct {
	Count int `json:"count"`
}

type requestCodesResponse struct {
	Codes   []string  `json:"codes"`
	Expires time.Time `json:"expires"`
}

func (app *Application) requestCodes(w http.ResponseWriter, r *http.Request) {
	log := getRequestLog(r)

	var request requestCodesRequest
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

	if request.Count < 1 || request.Count > 100 {
		log.WithField("count", request.Count).Error("Illegal number of codes requested")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	response := requestCodesResponse{
		Expires: time.Now().Add(config.ExpireHealthAuthorisationTokensAfter),
		Codes:   make([]string, request.Count),
	}
	for i := range response.Codes {
		response.Codes[i] = tools.GenerateCode()
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.WithField("error", err).Error("Error encoding response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
