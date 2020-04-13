package main

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
)

type contextKey string

const ctxLog = contextKey("requestLog")

func getRequestLog(r *http.Request) *logrus.Entry {
	log, ok := r.Context().Value(ctxLog).(*logrus.Entry)
	if ok {
		return log
	}
	return logrus.WithFields(logrus.Fields{
		"ip":  nil,
		"uri": nil,
	})
}

func (app *Application) addLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := app.log.WithFields(logrus.Fields{
			"ip":  r.RemoteAddr,
			"uri": r.RequestURI,
		})
		ctx := context.WithValue(r.Context(), ctxLog, log)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
