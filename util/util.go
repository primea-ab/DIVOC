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

func WriteBytes(w http.ResponseWriter, bytes []byte) {
	w.Header().Add("Content-Type", "application/octet-stream")

	w.Write(bytes)
}

func WriteError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)

	json.NewEncoder(w).Encode(map[string]string{
		"message": err.Error(),
	})
}

func QueryParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}
