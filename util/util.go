package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"divoc.primea.se/models"
)

var ServerAddress string
var ChunkByteSize int64 = 100
var RootShareFolder string = "share_folder"

func GetJSON(url string, v interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("got non-200 status code: %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(v)
}

func GetBytes(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}

	if resp.StatusCode != 200 {
		return []byte{}, fmt.Errorf("got non-200 status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

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

func RemoveStringFromSlice(value string, array []string) []string {
	returnArray := make([]string, 0)
	for _, arrayValue := range array {
		if arrayValue != value {
			returnArray = append(returnArray, arrayValue)
		}
	}
	return returnArray
}

func RemoveKeyFromMap(key string, m map[string]*models.SharedFile) map[string]*models.SharedFile {
	returnMap := make(map[string]*models.SharedFile, 0)
	for k, v := range m {
		if k != key {
			returnMap[k] = v
		}
	}
	return returnMap
}

func Keys(m map[string]struct{}) []string {
	keys := make([]string, 0)

	for k := range m {
		keys = append(keys, k)
	}

	return keys
}
