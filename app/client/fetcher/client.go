package fetcher

import (
	"fmt"
	"math/rand"
	"net/http"

	"divoc.primea.se/models"
)

type Client interface {
	GetShard(id int64, metaData models.ResultFile) ([]byte, error)
}

type ShardClient struct {
	rand *rand.Rand
}

func (s *ShardClient) GetShard(id int64, metaData models.ResultFile) ([]byte, error) {
	clientId := s.rand.Intn(len(metaData.Clients))

	url := fmt.Sprintf("%s:8080?chunk=%d&hash=%s", metaData.Clients[clientId], id, metaData.Hash)
	fmt.Printf("Make req to %+v\n", url)
	res, err := http.Get(url)

	if err != nil {
		fmt.Printf("Request failed with error %+v\n", err)
		return nil, err
	}

	var body []byte
	res.Body.Read(body)
	res.Body.Close()
	return body, nil
}

func NewShardClient() *ShardClient {
	return &ShardClient{rand: rand.New(rand.NewSource(1337))}
}
