package client

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"divoc.primea.se/app/client/fetcher"
	"divoc.primea.se/app/client/shardclient"
	"divoc.primea.se/models"
	"divoc.primea.se/util"
)

var (
	fileHashTable = make(map[string]string)
)

func StartClient() {
	registerContentOfFolder()

	http.HandleFunc("/download", downloadHandler)

	go func() {
		log.Fatal(http.ListenAndServe(":3001", nil))
	}()

	for {
		fmt.Print("\nSearch: ")
		var query string
		fmt.Scanln(&query)
		var searchResponse models.SearchResponse
		if err := util.GetJSON(util.ServerAddress+"/search?query="+query, &searchResponse); err != nil {
			log.Fatal(err)
		}
		for i, result := range searchResponse.Results {
			fmt.Printf("%d: %d %d [%s]\n", i, len(result.Clients), result.Size, strings.Join(result.Names, ", "))
		}
		fmt.Print("Download file: ")
		var fileIndex string
		fmt.Scanln(&fileIndex)
		fileIndexAsInt, err := strconv.ParseInt(fileIndex, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		fetcher := fetcher.New(shardclient.NewHttpClient().WithRetries(3), 4, &searchResponse.Results[fileIndexAsInt])
		fetcher.Download()
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

	data := make([]byte, util.ChunkByteSize)
	count, err := file.ReadAt(data, util.ChunkByteSize*chunkIndex)
	if err != nil && err != io.EOF {
		util.WriteError(w, err)
		return
	}

	util.WriteBytes(w, data[:count])
}

func registerContentOfFolder() {
	file, err := os.Open(util.RootShareFolder)
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	var metadataArray = []models.File{}

	fileInfos, _ := file.Readdir(0)
	for _, fileInfo := range fileInfos {
		hash := getHashForFile(util.RootShareFolder + "/" + fileInfo.Name())
		metadataArray = append(metadataArray, models.File{
			Name: fileInfo.Name(),
			Hash: hash,
			Size: fileInfo.Size(),
		})
		fileHashTable[hash] = util.RootShareFolder + "/" + fileInfo.Name()
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
