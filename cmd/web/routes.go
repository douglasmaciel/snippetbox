package main

import (
	"net/http"

	"snippetbox.douglasmaciel.com/ui"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))
	mux.Handle("GET /{$}", app.sessionManager.LoadAndSave(http.HandlerFunc(app.home)))
	mux.Handle("GET /snippets/{id}", app.sessionManager.LoadAndSave(http.HandlerFunc(app.snippetView)))
	mux.Handle("GET /snippet", app.sessionManager.LoadAndSave(http.HandlerFunc(app.snippetCreate)))
	mux.Handle("POST /snippet", app.sessionManager.LoadAndSave(http.HandlerFunc(app.snippetCreatePost)))
	return app.recoverPanic(app.logRequest(commonHeaders(mux)))
}
