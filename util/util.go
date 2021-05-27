package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var ServerAddress string

func PostJSON(path string, v interface{}) error {
	body, err := json.Marshal(v)
	if err != nil {
		return err
	}

	resp, err := http.Post(ServerAddress+path, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("got non-200 status code: %d", resp.StatusCode)
	}

	return nil
}

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

	WriteJSON(w, map[string]string{
		"message": err.Error(),
	})
}

func QueryParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}
