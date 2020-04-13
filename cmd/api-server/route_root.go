package main

import "net/http"

func (app *Application) root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	_, _ = w.Write([]byte("OK\n"))
}
