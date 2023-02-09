package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form creates a custome form struct, embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// Valid return true if there are no errors
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// New intializes a form struct
func New(data url.Values) *Form {
	return &Form{
		data, 
		errors(map[string][]string{}),
	}
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field is required")
		}
	}
}

// Has check if form field in post request is not empty
func (f *Form) Has(field string, req *http.Request) bool {
	x := req.Form.Get(field)
	if x == "" {
		f.Errors.Add(field, "This field cannot blank")
		return false
	} else {
		return true
	}
}

func (f *Form) MinLength(field string, length int, req *http.Request) bool {
	x := req.Form.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d character long", length))
		return false
	}
	return true
}

func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}