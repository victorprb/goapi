package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	mux := mux.NewRouter()
	mux.HandleFunc("/ping", ping).Methods("POST")
	mux.HandleFunc("/auth", app.auth).Methods("POST")
	mux.HandleFunc("/user", app.createUser).Methods("POST")
	s := mux.PathPrefix("/user").Subrouter()
	s.HandleFunc("/{uuid}", app.showUser).Methods("GET")
	s.Use(app.authenticateJWTMiddleware)
	
	//	mux.HandleFunc("/user", app.updateUser).Methods("PUT")
	//	mux.HandleFunc("/user", app.deleteUser).Methods("DELETE")
	//	mux.HandleFunc("/user/{:uuid}", app.showUser).Methods("GET")

	return standardMiddleware.Then(mux)
}
