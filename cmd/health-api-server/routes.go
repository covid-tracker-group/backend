package main

import (
	"net/http"
	"net/http/httputil"

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
	standardMiddleware := alice.New(app.addLog)
	apiMiddleware := alice.New(rateLimiter.RateLimit)

	router := mux.NewRouter()
	router.HandleFunc("/_status/healthy", app.healthy).Methods("GET")
	router.Handle("/api/request-codes", apiMiddleware.ThenFunc(app.requestCodes)).Methods("POST")

	if app.config.ProxyURL != nil {
		proxy := httputil.NewSingleHostReverseProxy(app.config.ProxyURL)
		router.HandleFunc("/{path:.*}", func(w http.ResponseWriter, r *http.Request) {
			r.URL.Path = mux.Vars(r)["path"]
			proxy.ServeHTTP(w, r)
		}).Methods("GET")
	} else {
		router.Handle("/", http.FileServer(http.Dir(app.config.HttpPath)))
	}

	return standardMiddleware.Then(router)
}
