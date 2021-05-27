package fetcher

import (
	"divoc.primea.se/models"
)

type Client interface {
	GetShard(hash string, id int64) ([]byte, error)
	GetFileMetaData(filename string) (models.FileMeta, error)
}
