package validator

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	Errors   map[string]string
	Warnings []string
}

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func (v *Validator) IsValid() bool {
	return len(v.Errors) == 0 && len(v.Warnings) == 0
}

func (v *Validator) AddError(key, msg string) {
	if v.Errors == nil {
		v.Errors = make(map[string]string)
	}

	v.Errors[key] = msg
}

func (v *Validator) AddWarning(msg string) {
	v.Warnings = append(v.Warnings, msg)
}

func (v *Validator) Validate(ok bool, key, msg string) {
	if !ok {
		v.AddError(key, msg)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MinChars(value string, limit int) bool {
	return utf8.RuneCountInString(value) >= limit
}

func MaxChars(value string, limit int) bool {
	return utf8.RuneCountInString(value) <= limit
}

func InList[T comparable](value T, list ...T) bool {
	return slices.Contains(list, value)
}

func MatchRegexp(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func IsEmail(value string) bool {
	return MatchRegexp(value, EmailRX)
}
