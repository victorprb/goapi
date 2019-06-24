package main

import (
	"bytes"
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.post(t, "POST", "/ping", "", nil)

	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}

	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}

func TestShowUser(t *testing.T) {
	// Create a new instance of our application struct which uses the mocked
	// dependencies.
	app := newTestApplication(t)
	token := newTestTokenJWT(t)

	// Establish a new test server for running end-to-end tests.
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Set up some table-driven tests to check the responses sent by our
	// application for different URLs.
	tests := []struct {
		name string
		urlPath string
		jwtToken string
		wantCode int
		wantBody []byte
	}{
		{"Valid UUID", "/user/29c31a03-f111-4f6f-b49e-07b0117fbb42", token, http.StatusOK, []byte("Cloud")},
		{"Non-existent ID", "/user/3aa9289e-d80f-4a67-bc77-2b68271fb98b", token, http.StatusNotFound, nil},
		{"Trailing slash", "/user/29c31a03-f111-4f6f-b49e-07b0117fbb42/", token, http.StatusNotFound, nil},
		{"Invalid token", "/user/29c31a03-f111-4f6f-b49e-07b0117fbb42", "invalid_token", http.StatusInternalServerError, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.post(t, "GET", tt.urlPath, tt.jwtToken, nil)
	
			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}
	
			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contains %q", tt.wantBody)
			}
		})
	}
}