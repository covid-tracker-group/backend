package main

import (
	"net/http"
)

func addSecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// We allow connect-src to support the webpack-dev-server with hot reload
		w.Header().Add("Content-Security-Policy", "default-src 'self'; connect-src 'self'; style-src 'self' 'unsafe-inline'; form-action 'none'; frame-ancestors 'none'; block-all-mixed-content;")
		w.Header().Add("Feature-Policy", "autoplay 'none'; camera 'none'; display-capture 'none'; document-domain 'none'; encrypted-media 'none'; fullscreen 'self'; geolocation 'none'; microphone 'none'; midi 'none'; payment 'none'; vr 'none'")
		w.Header().Add("X-Frame-Options", "DENY")
		w.Header().Add("X-Content-Type-Options", "nosniff")
		w.Header().Add("Referrer-Policy", "same-origin")
		next.ServeHTTP(w, r)
	})
}
