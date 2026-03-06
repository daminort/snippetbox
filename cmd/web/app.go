package main

import (
	"html/template"
	"log/slog"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"

	"snippetbox.demien.net/internal/models"
)

type Application struct {
	logger         *slog.Logger
	snippets       models.SnippetServiceInstance
	users          models.UsersServiceInstance
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func (app *Application) LogRequest(r *http.Request) {
	var (
		ip     = r.RemoteAddr
		proto  = r.Proto
		method = r.Method
		uri    = r.URL.RequestURI()
	)
	app.logger.Info("incoming request",
		"ip", ip,
		"method", method,
		"uri", uri,
		"proto", proto,
	)
}
