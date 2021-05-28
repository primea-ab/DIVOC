package shardclient

import (
	"fmt"
	"math/rand"

	"divoc.primea.se/models"
	"divoc.primea.se/util"
)

type httpClient struct {
	rand *rand.Rand
}

func (s *httpClient) sample(data []string) string {
	id := s.rand.Intn(len(data))
	return data[id]
}

func (s *httpClient) GetShard(id int64, metaData models.ResultFile) ([]byte, error) {
	seeder := s.sample(metaData.Clients)
	url := fmt.Sprintf("http://%s:3001/download?chunk=%d&hash=%s", seeder, id, metaData.Hash)
	return util.GetBytes(url)
}

func NewHttpClient() *httpClient {
	return &httpClient{rand: rand.New(rand.NewSource(1337))}
}
