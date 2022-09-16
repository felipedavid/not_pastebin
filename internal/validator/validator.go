package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// Regex for checking the format of an email address (recommended by W3C and
// Web Hypertext.
var EmailRegex = regexp.MustCompile("/^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/")

type Validator struct {
	FieldErrors map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// CheckField checks if the expression represented by "ok" is valid. If it is not,
// An error is created in the FieldErrors map
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func NotBlank(val string) bool {
	return strings.TrimSpace(val) != ""
}

// MinChars returns true if the string contains at least 'n' characters
func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// MaxChars returns true if the string contains less than 'n' characters
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func PermittedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}

func Matches(value string, re *regexp.Regexp) bool {
	return re.MatchString(value)
}
