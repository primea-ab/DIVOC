package fetcher

import (
	"fmt"
	"math/rand"

	"divoc.primea.se/models"
	"divoc.primea.se/util"
)

type Client interface {
	GetShard(id int64, metaData models.ResultFile) ([]byte, error)
}

type ShardClient struct {
	rand *rand.Rand
}

func (s *ShardClient) GetShard(id int64, metaData models.ResultFile) ([]byte, error) {
	clientId := s.rand.Intn(len(metaData.Clients))

	url := fmt.Sprintf("http://%s:3001/download?chunk=%d&hash=%s", metaData.Clients[clientId], id, metaData.Hash)
	fmt.Printf("Make req to %+v\n", url)

	return util.GetBytes(url)
}

func NewShardClient() *ShardClient {
	return &ShardClient{rand: rand.New(rand.NewSource(1337))}
}
