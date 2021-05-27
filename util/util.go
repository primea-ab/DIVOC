package util

import (
	"encoding/json"
	"net/http"
)

func ReadJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func WriteJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Add("Content-Type", "application/json")

	json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)

	json.NewEncoder(w).Encode(map[string]string{
		"message": err.Error(),
	})
}
