package main

import (
	"encoding/json"
	"net/http"

	appcore "github.com/mrbagir/qcash-appcore/pkg/app"
)

func main() {
	app := appcore.New()

	app.Handle("POST /api/sayhello", HalloHandler)

	app.Run()
}

func HalloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	response := map[string]string{"message": "Hello " + name}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
