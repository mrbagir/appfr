package main

import (
	"encoding/json"
	"net/http"

	appcore "github.com/mrbagir/appfr"
)

func main() {
	app := appcore.New()

	app.Handle("POST /api/sayhello", HelloHandler)

	app.Run()
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Hello " + name,
	})
}
