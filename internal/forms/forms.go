package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Form creates a custom form struct , embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// Valid returns true if there are no errors , otherwise false.
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// New init a form struct
func New(data url.Values) *Form {
	return &Form{data, errors(map[string][]string{})}
}

// Required checks for required field
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "Field "+field+" is required")
		}
	}
}

// Has checks if form filed is in post and not empty
func (f *Form) Has(field string, request *http.Request) bool {
	x := request.Form.Get(field)
	if x == "" {
		return false
	}
	return true
}

// MinLength checks for string minimum length
func (f *Form) MinLength(field string, length int, request *http.Request) bool {
	x := request.Form.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be atleast %d characters long.", length))
		return false
	}
	return true
}
