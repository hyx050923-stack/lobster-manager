package api

import (
	"encoding/json"
	"net/http"

	"ClawButler/internal/container"
)

func ListContainers(w http.ResponseWriter, r *http.Request) {

	manager, _ := container.NewManager()

	list, _ := manager.List()

	w.Write(list)
}