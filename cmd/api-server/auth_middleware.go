package main

import (
	"context"
	"net/http"
	"strings"
)

const ctxTracingAuthCode = contextKey("tracingAuthCode")

func getTracingAuthenticationCode(r *http.Request) string {
	authCode, _ := r.Context().Value(ctxTracingAuthCode).(string)
	return authCode
}

func (app *Application) requireTracingAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		var token string
		hasAuth := false
		if strings.ToLower(authHeader[:7]) == "bearer " {
			var err error
			token = authHeader[7:]
			hasAuth, err = app.tokenManager.VerifyToken(token)
			if err != nil {
				app.log.Errorf("Error verifying auth bearer: %v", err)
			}
		}
		if !hasAuth {
			app.log.WithField("Authorization", authHeader).Warn("Attempt to access endpoint without valid auth bearer")
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		} else {
			ctx := context.WithValue(r.Context(), ctxTracingAuthCode, token)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
