package shardclient

import "divoc.primea.se/models"

type Client interface {
	GetShard(id int64, metaData models.ResultFile) ([]byte, error)
}
