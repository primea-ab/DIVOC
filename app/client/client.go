package client

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"divoc.primea.se/models"
	"divoc.primea.se/util"
)

func StartClient() {
	registerContentOfFolder()

	http.HandleFunc("/download", downloadHandler)

	fmt.Println("Server listens on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	var chunkIndex int64 = 0

	file, err := os.Open("share_folder/file1.txt")
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	data := make([]byte, 100)
	count, err := file.ReadAt(data, 100*chunkIndex)
	if err != nil && err != io.EOF {
		util.WriteError(w, err)
		return
	}

	util.WriteBytes(w, data[:count])
}

func registerContentOfFolder() {
	file, err := os.Open("share_folder")
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	var metadataArray = []models.File{}

	fileInfos, _ := file.Readdir(0)
	for _, fileInfo := range fileInfos {
		metadataArray = append(metadataArray, models.File{
			Name: fileInfo.Name(),
			Hash: getHashForFile("share_folder/" + fileInfo.Name()),
			Size: fileInfo.Size(),
		})
	}

	var request = models.RegisterRequest{
		Files: metadataArray,
	}

	body, err := json.Marshal(request)
	if err != nil {
		fmt.Println(err)
	}
	http.Post("http://192.168.1.235:8080/register", "application/json", bytes.NewBuffer(body))
}

func getHashForFile(filepath string) string {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
	}

	data := make([]byte, fileInfo.Size())
	count, err := file.Read(data)
	if err != nil {
		log.Fatal(err)
	}

	first := sha256.New()
	first.Write([]byte(data[:count]))

	return fmt.Sprintf("%x", first.Sum(nil))
}
