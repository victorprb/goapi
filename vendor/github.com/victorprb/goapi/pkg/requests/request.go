package requests

import (
	"reflect"
	"regexp"
	"strings"
)

// Email validate regex pattern
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Request struct {
	BodyFields interface{}
	Errors     errors
}

// New initializes a Request struct with data
func New(data interface{}) *Request {
	return &Request{
		data,
		errors(map[string][]string{}),
	}
}

// Required check fields in request body
func (r *Request) Required(fields ...string) {
	for _, field := range fields {
		st := reflect.ValueOf(r.BodyFields)
		value := st.FieldByName(field).String()
		if strings.TrimSpace(value) == "" {
			r.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// MatchesPattern check if a string matchs a regex pattern
func (r *Request) MatchesPattern(field string, pattern *regexp.Regexp) {
	st := reflect.ValueOf(r.BodyFields)
	value := st.FieldByName(field).String()
	if strings.TrimSpace(value) == "" {
		return
	}
	if !pattern.MatchString(value) {
		r.Errors.Add(field, "This field is invalid")
	}
}

// Valid check if there's no errors on validations
func (r *Request) Valid() bool {
	return len(r.Errors) == 0
}
