package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetbox.demien.net/internal/models"
)

// Snippets

func (app *Application) snippetsList(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.GetLatest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, r, http.StatusOK, "home.html", data)
}

func (app *Application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	snippet, err := app.snippets.GetSnippet(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.NotFound(w, r)
			return
		}

		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, r, http.StatusOK, "view.html", data)
}

func (app *Application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = snippetForm{
		Expires: 365,
	}
	app.render(w, r, http.StatusOK, "create.html", data)
}

func (app *Application) snippetInsert(w http.ResponseWriter, r *http.Request) {
	var form snippetForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form = validateSnippetForm(form)
	if len(form.Errors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	snippet, err := app.snippets.InsertSnippet(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", snippet.ID), http.StatusSeeOther)
}

// Users

func (app *Application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = signupForm{}
	app.render(w, r, http.StatusOK, "signup.html", data)
}

func (app *Application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = loginForm{}
	app.render(w, r, http.StatusOK, "login.html", data)
}

func (app *Application) userInsert(w http.ResponseWriter, r *http.Request) {
	var form signupForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form = validateSignupForm(form)
	if len(form.Errors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}

	_, err = app.users.InsertUser(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors["email"] = "Email is already in use"

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "signup.html", data)
			return
		}

		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "User successfully created")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *Application) userVerify(w http.ResponseWriter, r *http.Request) {
	var form loginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	user, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Warnings = []string{"Email or Password is incorrect"}

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusOK, "login.html", data)
			return
		}

		app.serverError(w, r, err)
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "userID", user.ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) userLogout(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "userID")
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Health

func (app *Application) ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
