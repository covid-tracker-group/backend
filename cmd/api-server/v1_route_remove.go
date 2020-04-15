package main

import (
	"net/http"
)

func (app *Application) remove(w http.ResponseWriter, r *http.Request) {
	tracingAuthCode := getTracingAuthenticationCode(r)
	app.keyStorage.PurgeRecords(tracingAuthCode)
	app.tracingAuthTokenManager.RetractToken(tracingAuthCode)
	w.WriteHeader(http.StatusNoContent)
}
