package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"snippetbox.demien.net/internal/models"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	addr := flag.String("addr", ":8080", "http service address")

	// logger
	loggerHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})
	logger := slog.New(loggerHandler)

	// database
	db, err := openDB()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	// template cache
	templateCache, err := newTemplatesCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// form decoder
	formDecoder := form.NewDecoder()

	// session manager
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	// app initialization
	app := &Application{
		logger:         logger,
		snippets:       models.NewSnippetsService(db),
		users:          models.NewUsersService(db),
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// TLS config
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:       tls.VersionTLS13,
	}

	// HTTP server
	srv := &http.Server{
		Addr:         *addr,
		Handler:      app.createRoutes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.logger.Info("Listening on " + *addr)

	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	app.logger.Error(err.Error())
	os.Exit(1)
}

func openDB() (*sql.DB, error) {
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPwd := os.Getenv("DB_PASSWORD")
	dsn := fmt.Sprintf("%s:%s@/%s?parseTime=true", dbUser, dbPwd, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
