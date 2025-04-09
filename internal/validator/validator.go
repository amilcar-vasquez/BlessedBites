// file: internal/validator/validator.go
package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// validator type is just a map
type Validator struct {
	Errors map[string]string
}

// use factory function for the validator
func NewValidator() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

// check if data is good
func (v *Validator) ValidData() bool {
	return len(v.Errors) == 0
}

// add an error entry to the error map
func (v *Validator) AddError(field string, message string) {
	if _, exists := v.Errors[field]; !exists {
		v.Errors[field] = message
	}
}

// add an error if the validation check fails
func (v *Validator) Check(ok bool, field string, message string) {
	if !ok {
		v.AddError(field, message)
	}
}

// returns true if data is present in the input box
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxLength returns true if the value contains no mor than n characters
func MaxLength(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// MinLength returns true if the value contains at least n characters
func MinLength(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// IsEmail returns true if the value is a valid email address
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func IsEmail(value string) bool {
	return EmailRX.MatchString(value)
}
