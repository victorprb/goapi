package requests

import (
	"reflect"
	"regexp"
	"strings"
)

var EmaiRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Request struct {
	BodyFields interface{}
	Errors     errors
}

func New(data interface{}) *Request {
	return &Request{
		data,
		errors(map[string][]string{}),
	}
}

func (r *Request) Required(fields ...string) {
	for _, field := range fields {
		st := reflect.ValueOf(r.BodyFields)
		value := st.FieldByName(field).String()
		if strings.TrimSpace(value) == "" {
			r.Errors.Add(field, "This field cannot be blank")
		}
	}
}

func (r *Request) Valid() bool {
	return len(r.Errors) == 0
}
