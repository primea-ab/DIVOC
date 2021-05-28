package client

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"divoc.primea.se/app/client/fetcher"
	"divoc.primea.se/app/client/shardclient"
	"divoc.primea.se/models"
	"divoc.primea.se/util"
)

var (
	fileHashTable = make(map[string]string)
	aliveTimer    *time.Timer
	progressMap   = make(map[string]*fetcher.FileFetcher)
)

func StartClient() {
	http.HandleFunc("/alive", alive)
	http.HandleFunc("/download", download)
	http.HandleFunc("/search", search)
	http.HandleFunc("/start-download", startDownload)
	http.HandleFunc("/progress", progressHandler)

	http.Handle("/", http.FileServer(http.Dir("app/client/static")))

	time.AfterFunc(time.Second, func() {
		aliveTimer = time.NewTimer(5 * time.Second)
		go ensureServerConnection()

		registerContentOfFolder()
	})

	log.Fatal(http.ListenAndServe(":3001", nil))
}

func startDownload(w http.ResponseWriter, r *http.Request) {
	var resultFile models.ResultFile
	if err := util.ReadJSON(r, &resultFile); err != nil {
		util.WriteError(w, err)
		return
	}

	fetcher := fetcher.New(shardclient.NewHttpClient().WithRetries(3), 4, &resultFile)
	progressMap[resultFile.Names[0]] = fetcher

	fetcher.Download()
}

func search(w http.ResponseWriter, r *http.Request) {
	var searchResponse models.SearchResponse
	if err := util.GetJSON(util.ServerAddress+"/search?query="+util.QueryParam(r, "query"), &searchResponse); err != nil {
		log.Fatal(err)
	}

	util.WriteJSON(w, searchResponse)
}

func alive(w http.ResponseWriter, r *http.Request) {
	aliveTimer.Stop()
	<-r.Context().Done()
	aliveTimer.Reset(5 * time.Second)
}

func ensureServerConnection() {
	<-aliveTimer.C
	registerContentOfFolder()
}

func download(w http.ResponseWriter, r *http.Request) {
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

type progressItem struct {
	Name     string
	Progress float64
}

func progressHandler(w http.ResponseWriter, r *http.Request) {
	progressItems := make([]progressItem, 0)
	for name, fetcher := range progressMap {
		progressItems = append(progressItems, progressItem{
			Name:     name,
			Progress: fetcher.Progress,
		})
	}

	util.WriteJSON(w, progressItems)
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

	if err := util.PostJSON("/register", request); err != nil {
		log.Fatal("register to server failed")
	}
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
