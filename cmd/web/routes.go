package main

import (
	"net/http"

	"github.com/justinas/alice"

	"snippetbox.demien.net/ui"
)

func (app *Application) createRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	dynamic := alice.New(app.sessionManager.LoadAndSave, app.preventCSRF, app.authenticate)
	protected := dynamic.Append(app.requireAuth)

	// Snippets
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.snippetsList))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /snippet/create", protected.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetInsert))

	// Users
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userInsert))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userVerify))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogout))

	// Health
	mux.HandleFunc("GET /ping", app.ping)

	standard := alice.New(app.recoverPanic, app.requestLogger, commonHeaders)

	return standard.Then(mux)
}
