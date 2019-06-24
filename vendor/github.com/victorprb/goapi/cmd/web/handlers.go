package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/victorprb/goapi/pkg/models"
	"github.com/victorprb/goapi/pkg/requests"

	"github.com/gorilla/mux"
	jwt "github.com/dgrijalva/jwt-go"
)

var response map[string]interface{}

// Response to check if the server is up
func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

// Authenticate the user and return a JWT token
func (app *application) auth(w http.ResponseWriter, r *http.Request) {
	var cred *models.Credentials
	
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = json.Unmarshal(body, &cred)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	request := requests.New(*cred)
	// Check the required fields of body request
	request.Required("Login", "Password")

	if !request.Valid() {
		app.writeJSON(w, http.StatusBadRequest, request.Errors)
		return
	}

	// Check the user credentials in DB
	uuid, err := app.users.Authenticate(cred)
	if err == models.ErrInvalidCredentials {
		// request.Errors.Add("generic", "Email or password is incorrect")
		response = map[string]interface{}{
			"status": "error",
			"message": "User couldn't authenticate",
		}
		app.writeJSON(w, http.StatusUnauthorized, response)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	// Define token expiration time
	expirationTime := time.Now().Add(time.Hour * 24)
	// Create the JWT claims, which includes the username and expiry time
	claims := &models.Claims{
		Login: cred.Login,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the alg user for signing and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	tokenString, err := token.SignedString([]byte(app.secretKey))
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Get the user info from DB
	user, err := app.users.Get(uuid)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	
	// Response success with the new JWT token info
	response = map[string]interface{}{
		"status": "success",
		"message": "User found and token generated",
		"tokenjwt": tokenString,
		"expires": expirationTime.Format(time.RFC3339),
		"login": user,
	}
	app.writeJSON(w, http.StatusOK, response)
}

// Create a new user
func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
	var user *models.User

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	request := requests.New(*user)
	// Check the required fields of body request
	request.Required("FirstName", "LastName", "CPF", "Email")
	request.MatchesPattern("Email", requests.EmailRX)

	if !request.Valid() {
		app.writeJSON(w, http.StatusBadRequest, request.Errors)
		return
	}

	// Create the new user in DB
	err = app.users.Insert(user)
	// Check if email or cpf already exists
	if err == models.ErrDuplicateEmail {
		request.Errors.Add("email", "Address is already in use")
		app.writeJSON(w, http.StatusInternalServerError, request.Errors)
		return
	} else if err == models.ErrDuplicateCPF {
		request.Errors.Add("cpf", "Number is already in use")
		app.writeJSON(w, http.StatusInternalServerError, request.Errors)
		return
	} else if err != nil {
		app.serverError(w, err)
	}

	app.writeJSON(w, http.StatusCreated, fmt.Sprintf("UUID: %v", user.UUID))
}

// Retrieve user info
func (app *application) showUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]
	if uuid == "" {
		app.notFound(w)
		return
	}

	// Get the user info from DB
	u, err := app.users.Get(uuid)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	
	app.writeJSON(w, http.StatusOK, u)
}
