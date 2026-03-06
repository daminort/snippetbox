package main

import (
	"snippetbox.demien.net/internal/validator"
)

func validateSnippetForm(form snippetForm) snippetForm {
	v := validator.Validator{}

	v.Validate(validator.NotBlank(form.Title), "title", "Title cannot be blank")
	v.Validate(validator.MaxChars(form.Title, 100), "title", "Title must be less than 100 characters")
	v.Validate(validator.NotBlank(form.Content), "content", "Content cannot be blank")
	v.Validate(validator.InList(form.Expires, 1, 7, 365), "expires", "Expires should be either 365 or 7 or 1")

	form.Errors = make(map[string]string)
	if !v.IsValid() {
		form.Errors = v.Errors
	}

	return form
}

func validateSignupForm(form signupForm) signupForm {
	v := validator.Validator{}

	v.Validate(validator.NotBlank(form.Name), "name", "Name cannot be blank")
	v.Validate(validator.NotBlank(form.Email), "email", "Email cannot be blank")
	v.Validate(validator.NotBlank(form.Password), "password", "Password cannot be blank")

	v.Validate(validator.IsEmail(form.Email), "email", "Email should be a valid email address")
	v.Validate(validator.MinChars(form.Password, 8), "password", "Password should be at least 8 characters")

	form.Errors = make(map[string]string)
	if !v.IsValid() {
		form.Errors = v.Errors
	}

	return form
}

func validateLoginForm(form loginForm) loginForm {
	v := validator.Validator{}

	v.Validate(validator.NotBlank(form.Email), "email", "Email cannot be blank")
	v.Validate(validator.NotBlank(form.Password), "password", "Password cannot be blank")
	v.Validate(validator.IsEmail(form.Email), "email", "Email should be a valid email address")

	form.Errors = make(map[string]string)
	if !v.IsValid() {
		form.Errors = v.Errors
	}

	return form
}
