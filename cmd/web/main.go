package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"snippetbox.douglasmaciel.com/internal/models"
)

type application struct {
	logger        *slog.Logger
	snippets      *models.SnippetModel
	templateCache templateCache
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	app := &application{
		logger:        logger,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}
	logger.Info("Starting server", "addr", *addr)
	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
