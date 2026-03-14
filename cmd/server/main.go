package main

import (
	"log"
	"net/http"

	"ClawButler/internal/api"
)

func main() {

	router := api.NewRouter()

	log.Println("server started")

	http.ListenAndServe("127.0.0.1:8080", router)
}