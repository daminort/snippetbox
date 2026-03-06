package main

import (
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"time"

	"github.com/justinas/nosurf"
	"snippetbox.demien.net/internal/models"
	"snippetbox.demien.net/ui"
)

type templateData struct {
	Year            int
	Snippet         *models.Snippet
	Snippets        []models.Snippet
	User            *models.User
	IsAuthenticated bool
	Form            any
	Flash           string
	CSRFToken       string
}

type snippetForm struct {
	Title   string            `form:"title"`
	Content string            `form:"content"`
	Expires int               `form:"expires"`
	Errors  map[string]string `form:"-"`
}

type signupForm struct {
	Name     string            `form:"name"`
	Email    string            `form:"email"`
	Password string            `form:"password"`
	Errors   map[string]string `form:"-"`
}

type loginForm struct {
	Email    string            `form:"email"`
	Password string            `form:"password"`
	Errors   map[string]string `form:"-"`
	Warnings []string          `form:"-"`
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006 at 15:04")
}

func (app *Application) newTemplateData(r *http.Request) *templateData {
	user, ok := r.Context().Value(userKey).(*models.User)

	return &templateData{
		Year:            time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		User:            user,
		IsAuthenticated: ok,
		CSRFToken:       nosurf.Token(r),
	}
}

func newTemplatesCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.html",
			"html/partials/*.html",
			page,
		}

		tmpl, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = tmpl
	}

	return cache, nil
}
