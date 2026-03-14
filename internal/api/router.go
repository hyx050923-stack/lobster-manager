package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {

	r := mux.NewRouter()

	r.HandleFunc("/containers", ListContainers).Methods("GET")

	r.HandleFunc("/containers/run", RunContainer).Methods("POST")

	return r
}