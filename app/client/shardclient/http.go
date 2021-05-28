package shardclient

import (
	"fmt"
	"math/rand"

	"divoc.primea.se/models"
	"divoc.primea.se/util"
)

type httpClient struct {
	rand       *rand.Rand
	numRetries uint
}

func (s *httpClient) sample(data []string) string {
	id := s.rand.Intn(len(data))
	return data[id]
}

func (s *httpClient) GetShard(id int64, metaData models.ResultFile) ([]byte, error) {
	var dataFunc func(uint) ([]byte, error)

	dataFunc = func(retriesLeft uint) ([]byte, error) {
		seeder := s.sample(metaData.Clients)
		url := fmt.Sprintf("http://%s:3001/download?chunk=%d&hash=%s", seeder, id, metaData.Hash)
		res, err := util.GetBytes(url)
		if err != nil && retriesLeft < 1 {
			return nil, err
		}
		if err != nil {
			metaData.Clients = util.RemoveStringFromSlice(seeder, metaData.Clients)
			if len(metaData.Clients) == 0 {
				return nil, err
			}
			return dataFunc(retriesLeft - 1)
		}
		return res, nil
	}

	return dataFunc(s.numRetries)
}

func NewHttpClient() *httpClient {
	return &httpClient{rand: rand.New(rand.NewSource(1337)), numRetries: 1}
}

func (s *httpClient) WithRetries(num uint) *httpClient {
	s.numRetries = num
	return s
}
