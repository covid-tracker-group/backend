package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/throttled/throttled"
	"github.com/throttled/throttled/store/memstore"
)

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
	standardMiddleware := alice.New(app.addLog, rateLimiter.RateLimit)
	authMiddleware := alice.New(app.requireTracingAuthentication)

	router := mux.NewRouter()

	router.HandleFunc("/", app.root).Methods("GET")

	v1router := router.PathPrefix("/v1").Subrouter()
	v1router.HandleFunc("/report", app.report).Methods("POST")
	v1router.Handle("/submit-keys", authMiddleware.Then(http.HandlerFunc(app.submit))).Methods("POST")
	v1router.Handle("/remove", authMiddleware.Then(http.HandlerFunc(app.remove))).Methods("POST")

	return standardMiddleware.Then(router)
}
