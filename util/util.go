package util

import (
	"encoding/json"
	"net/http"
)

func WriteError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)

	json.NewEncoder(w).Encode(map[string]string{
		"message": err.Error(),
	})
}
