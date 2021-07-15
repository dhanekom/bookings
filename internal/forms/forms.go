package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form contains form errors and embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// Valid returns true if the form is valid, otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// New initializes a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors{}, // errors(map[string][]string{})
	}
}

// Required checks for required fields
func (f *Form) Required(fields ...string) bool {
	result := true
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
			result = false
		}
	}

	return result
}

// MinLength checks for string minimum length
func (f *Form) MinLength(field string, length int) bool {
	value := f.Get(field)
	if len(value) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

// IsEmail checks for valid email address
func (f *Form) IsEmail(field string) bool {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
		return false
	}
	return true
}

// Has check if form field is in post an not empty
func (f *Form) Has(field string) bool {
	x := f.Get(field)
	if x == "" {
		f.Errors.Add(field, "This field cannot be blank")
		return false
	}

	return true
}
