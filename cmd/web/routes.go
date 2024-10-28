package main

import (
	"net/http"
	"path/filepath"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.Handle("GET /{$}", app.sessionManager.LoadAndSave(http.HandlerFunc(app.home)))
	mux.Handle("GET /snippets/{id}", app.sessionManager.LoadAndSave(http.HandlerFunc(app.snippetView)))
	mux.Handle("GET /snippet", app.sessionManager.LoadAndSave(http.HandlerFunc(app.snippetCreate)))
	mux.Handle("POST /snippet", app.sessionManager.LoadAndSave(http.HandlerFunc(app.snippetCreatePost)))
	return app.recoverPanic(app.logRequest(commonHeaders(mux)))
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}
	s, _ := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}
	return f, nil
}
