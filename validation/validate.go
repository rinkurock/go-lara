package validate

import (
	"net/http"
	"strings"
)

type ValidationError struct {
	StatusCode int
	Status     bool
	Response   map[string]interface{}
}

func NewValidationError() ValidationError {
	v := ValidationError{}
	v.StatusCode = http.StatusUnprocessableEntity
	v.Status = false
	return v
}

func NewPasswordValidationError() ValidationError {
	v := ValidationError{}
	v.StatusCode = http.StatusConflict
	v.Status = false
	return v
}

func PasswordValidation(username string, password string) []string {
	errors := []string{}
	l := len(password)
	if username == password {
		errors = append(errors, "Phone & password can not be same")
	}
	if len(password) < 6 {
		errors = append(errors, "Must be at least 6 characters")
	}
	if len(password) > 128 {
		errors = append(errors, "Must be less than 128 characters")
	}
	if l != len(strings.TrimSpace(password)) {
		errors = append(errors, "Starting or ending with space are not allowed")
	}
	return errors
}
func UsernameValidation(username string) []string {
	errors := []string{}

	if !StringValidation(username) {
		errors = append(errors, "Valid username required")
	}
	return errors
}

func CountryIdValidation(countryId int) bool {
	if countryId != 0 {
		return true
	}
	return false
}

func StringValidation(string string) bool {
	n := strings.TrimSpace(string)
	if len(n) > 0 {
		return true
	}
	return false
}
