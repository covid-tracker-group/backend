package main

import "net/http"

func (app *Application) healthy(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("OK\n"))
}
