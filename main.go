package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sajid/handlers"
)

func main() {
	// new router
	r := mux.NewRouter()

	// arranging routes
	r.HandleFunc("/v1/users/{id}", handlers.GetUser).Methods("GET")
	r.HandleFunc("/v1/users", handlers.GetAllUsers).Methods("GET")
	r.HandleFunc("/v1/users", handlers.CreateUser).Methods("POST")
	r.HandleFunc("/v1/users/{id}", handlers.UpdateUser).Methods("PUT")
	r.HandleFunc("/v1/users/{id}", handlers.DeleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8081", r))
}
