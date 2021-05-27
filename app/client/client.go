package client

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"divoc.primea.se/models"
	"divoc.primea.se/util"
)

var (
	fileHashTable = make(map[string]string)
)

var serverAddress *string

func StartClient() {
	flag.StringVar(&util.ServerAddress, "server", "", "address to server")
	flag.Parse()

	//registerContentOfFolder()

	http.HandleFunc("/download", downloadHandler)

	go func() {
		log.Fatal(http.ListenAndServe(":3001", nil))
	}()

	for {
		fmt.Print("Query: ")
		var query string
		fmt.Scanln(&query)
		var searchResponse models.SearchResponse
		if err := util.GetJSON(util.ServerAddress+"/search?query="+query, &searchResponse); err != nil {
			log.Fatal(err)
		}
		for i, result := range searchResponse.Results {
			fmt.Printf("%d: %d %d [%s]", i, len(result.Clients), result.Size, strings.Join(result.Names, ", "))
		}
		fmt.Print("File index: ")
		var fileIndex string
		fmt.Scanln(&fileIndex)
	}
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	chunkIndex, err := strconv.ParseInt(util.QueryParam(r, "chunk"), 0, 64)
	hash := util.QueryParam(r, "hash")

	if err != nil {
		util.WriteError(w, err)
		return
	}

	file, err := os.Open(fileHashTable[hash])
	if err != nil {
		util.WriteError(w, err)
		return
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
		hash := getHashForFile("share_folder/" + fileInfo.Name())
		metadataArray = append(metadataArray, models.File{
			Name: fileInfo.Name(),
			Hash: hash,
			Size: fileInfo.Size(),
		})
		fileHashTable[hash] = "share_folder/" + fileInfo.Name()
	}

	var request = models.RegisterRequest{
		Files: metadataArray,
	}

	util.PostJSON("/register", request)
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
