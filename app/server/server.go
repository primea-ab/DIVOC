package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"divoc.primea.se/models"
	"divoc.primea.se/util"
)

var (
	sharedFiles = make(map[string]*models.SharedFile)
	// Map with ip adresses as key and file hashes as values
	hostFiles = make(map[string][]string)
)

func StartServer() {
	http.HandleFunc("/register", register)
	http.HandleFunc("/search", search)

	log.Fatal(http.ListenAndServe(":3000", nil))
}

func register(w http.ResponseWriter, r *http.Request) {
	var registerRequest models.RegisterRequest
	if err := util.ReadJSON(r, &registerRequest); err != nil {
		util.WriteError(w, err)
		return
	}

	ipParts := strings.Split(r.RemoteAddr, ":")
	ip := ipParts[0]

	removeSeederFromFile(ip)

	for _, file := range registerRequest.Files {
		sharedFile, ok := sharedFiles[file.Hash]
		if !ok {
			sharedFile = &models.SharedFile{
				Size:  file.Size,
				Names: make(map[string]struct{}),
			}

			sharedFiles[file.Hash] = sharedFile
			_, ok := hostFiles[ip]
			if !ok {
				hostFiles[ip] = []string{file.Hash}
			} else {
				hostFiles[ip] = append(hostFiles[ip], file.Hash)
			}
		}

		sharedFile.Names[file.Name] = struct{}{}
		sharedFile.Clients = append(sharedFile.Clients, ip)
	}

	go checkAlive(ip)
}

func checkAlive(ip string) {
	url := fmt.Sprintf("http://%s:3001/alive", ip)

	wasJustDead := false

	for {
		if _, err := http.Get(url); err != nil {
			if wasJustDead {
				fmt.Printf("Client with IP %s was disconnected\n", ip)
				removeSeederFromFile(ip)
				break
			} else {
				wasJustDead = true
			}
		} else {
			wasJustDead = false
		}
	}
}

func search(w http.ResponseWriter, r *http.Request) {
	searchResponse := models.SearchResponse{
		Results: make([]models.ResultFile, 0),
	}

	for hash, sharedFile := range sharedFiles {
		names := util.Keys(sharedFile.Names)

		for _, name := range names {
			if strings.Contains(strings.ToLower(name), strings.ToLower(util.QueryParam(r, "query"))) {
				searchResponse.Results = append(searchResponse.Results, models.ResultFile{
					Hash:    hash,
					Names:   names,
					Clients: sharedFile.Clients,
					Size:    sharedFile.Size,
				})

				break
			}
		}
	}

	util.WriteJSON(w, searchResponse)
}

func removeSeederFromFile(ip string) {
	fileHashes, ok := hostFiles[ip]
	if !ok {
		return
	}
	for _, fileHash := range fileHashes {
		sharedFiles[fileHash].Clients = util.RemoveStringFromSlice(ip, sharedFiles[fileHash].Clients)
		if len(sharedFiles[fileHash].Clients) == 0 {
			sharedFiles = util.RemoveKeyFromMap(fileHash, sharedFiles)
		}

		hostFiles[ip] = util.RemoveStringFromSlice(fileHash, hostFiles[ip])
	}
}
