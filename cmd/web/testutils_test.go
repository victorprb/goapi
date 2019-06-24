package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	// "time"

	"github.com/victorprb/goapi/pkg/models/mock"

	jwt "github.com/dgrijalva/jwt-go"
)

var secretKey = "my_test_secret_key"

// Create a newTestApplication helper which returns an instance of our
// application struct containing mocked dependencies.
func newTestApplication(t *testing.T) *application {
	return &application{
		errorLog:      log.New(ioutil.Discard, "", 0),
		infoLog:       log.New(ioutil.Discard, "", 0),
		secretKey: 	   secretKey,
		users:         &mock.UserModel{},
	}
}

func newTestTokenJWT(t *testing.T) string {
	// expirationTime := time.Now().Add(time.Hour * 24)

	// Declare the token with the alg user for signing and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{})

	// Create the JWT string
	tokenString, _ := token.SignedString([]byte(secretKey))
	return tokenString
}

// Define a custom testServer type which anonymously embeds a httptest.Server
// instance.
type testServer struct {
	*httptest.Server
}

// Create a newTestServer helper which initalizes and returns a new instance
// of our custom testServer type.
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewServer(h)

	// Disable redirect-following for the client. Essentially this function
	// is called after a 3xx response is received by the client, and returning
	// the http.ErrUseLastResponse error forces it to immediately return the
	// received response.
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

// Implement a get method on our custom testServer type. This makes a GET
// request to a given url path on the test server, and returns the response
// status code, headers and body.
func (ts *testServer) post(t *testing.T, method, urlPath, token string, payload []byte) (int, http.Header, []byte) {
	rq, err := http.NewRequest(method, ts.URL + urlPath, bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}
	rq.Header.Set("Content-Type", "application/json")
	rq.Header.Set("Authorization", "bearer " + token)
	rs, err := ts.Client().Do(rq)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, body
}
