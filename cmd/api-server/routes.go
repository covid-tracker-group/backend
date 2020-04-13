package main

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/throttled/throttled"
	"github.com/throttled/throttled/store/memstore"
)

func stripV1Prefix(h http.Handler) http.Handler {
	return http.StripPrefix("/v1", h)
}

func (app *Application) createRateLimiter() (*throttled.HTTPRateLimiter, error) {
	store, err := memstore.New(65536)
	if err != nil {
		return nil, err
	}
	quota := throttled.RateQuota{
		MaxRate:  throttled.PerMin(20),
		MaxBurst: 5,
	}
	rateLimiter, err := throttled.NewGCRARateLimiter(store, quota)
	if err != nil {
		return nil, err
	}
	return &throttled.HTTPRateLimiter{
		RateLimiter: rateLimiter,
		VaryBy: &throttled.VaryBy{
			RemoteAddr: true,
		},
	}, nil
}

func (app *Application) routes() http.Handler {
	rateLimiter, err := app.createRateLimiter()
	if err != nil {
		app.log.Fatalf("Error creating rate limited: %v", err)
	}
	standardMiddleware := alice.New(app.addLog, rateLimiter.RateLimit, stripV1Prefix)
	authMiddleware := alice.New(app.requireTracingAuthentication)

	mux := http.NewServeMux()

	mux.HandleFunc("/", app.root)
	mux.HandleFunc("/report", app.report)
	mux.Handle("/submit", authMiddleware.Then(http.HandlerFunc(app.submit)))

	return standardMiddleware.Then(mux)
}
