package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/victorprb/goapi/pkg/models"
	"github.com/victorprb/goapi/pkg/requests"
)

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

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
	request.Required("FirstName", "LastName", "CPF", "Email")

	if !request.Valid() {
		app.writeJSON(w, http.StatusBadRequest, request.Errors)
		return
	}

	err = app.users.Insert(user)
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
