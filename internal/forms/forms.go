package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator/v12"
)

// Form creates a custom form struct and embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// New initializes a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
		// errors(make(errors, 0)),
	}
}

// Required checks for required fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		fieldValue := f.Get(field)
		if strings.TrimSpace(fieldValue) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// Has checks if form field is in post and not empty
func (f *Form) Has(field string, r *http.Request) bool {
	fieldValue := r.Form.Get(field)
	if fieldValue == "" {
		return false
	}
	return true
}

// Valid returns true if there are no errors, otherwise false
func (f *Form) IsValid() bool {
	return len(f.Errors) == 0
}

// MinLength checks for string minimum length
func (f *Form) MinLength(field string, length int) bool {
	valueInField := f.Get(field)
	if len(valueInField) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

// IsEmailValidator checks for valid email address
func (f *Form) IsEmailValidator(field string) {
	if !govalidator.IsEmail(f.Get("email")) {
		f.Errors.Add(field, "Invalid email address")
	}
}
