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
		token := r.Header.Get("Authorization")
		hasAuth := false
		if strings.HasPrefix(token, "bearer ") {
			var err error
			token := token[7:]
			hasAuth, err = app.tokenManager.VerifyToken(token)
			if err != nil {
				app.log.Errorf("Error verifying auth bearer: %v", err)
			}
		}
		if !hasAuth {
			app.log.Warn("Attempt to access endpoint without valid auth bearer")
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		} else {
			ctx := context.WithValue(r.Context(), ctxTracingAuthCode, token)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
