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

	for _, file := range registerRequest.Files {
		sharedFile, ok := sharedFiles[file.Hash]
		if !ok {
			sharedFile = &models.SharedFile{
				Size:  file.Size,
				Names: make(map[string]struct{}),
			}

			sharedFiles[file.Hash] = sharedFile
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
		names := keys(sharedFile.Names)

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

func keys(m map[string]struct{}) []string {
	keys := make([]string, 0)

	for k := range m {
		keys = append(keys, k)
	}

	return keys
}
