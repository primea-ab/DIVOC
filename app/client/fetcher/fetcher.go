package fetcher

import (
	"errors"
	"log"
	"math"
	"os"
	"sync"

	"divoc.primea.se/models"
)

type FileFetcher struct {
	client     Client
	numWorkers int64
	shardLen   int64
	meta       *models.ResultFile
}

func handleError(err error) {
	if err != nil {
		log.Fatalf("Something went wrong: %v\n", err)
	}
}

func (a *FileFetcher) Download() error {
	file, err := a.getFile(int(a.meta.Size), a.meta.Names[0])
	handleError(err)
	defer file.Close()

	numShards := int64(math.Ceil(float64(a.meta.Size) / float64(a.shardLen)))

	shardChan := make(chan int64, numShards)
	dataChan := make(chan partialResult, a.numWorkers)
	var wg sync.WaitGroup
	wg.Add(int(numShards))

	var i int64
	for i = 0; i < numShards; i += 1 {
		shardChan <- i
	}

	go a.startWriteWorker(dataChan, file, &wg, a.shardLen)

	for i = 0; i < a.numWorkers; i += 1 {
		go a.startFetchWorker(dataChan, shardChan, i)
	}

	wg.Wait()

	return nil
}

func (a *FileFetcher) getFile(numBytes int, filename string) (*os.File, error) {
	file, err := os.Open(filename)
	if err != nil {
		return a.createEmpty(numBytes, filename)
	}

	info, _ := file.Stat()
	if info.Size() == int64(numBytes) {
		return file, nil
	}

	return a.createEmpty(numBytes, filename)
}

func (a *FileFetcher) createEmpty(numBytes int, filename string) (*os.File, error) {
	file, err := os.Create(filename)
	handleError(err)
	bytes := make([]byte, numBytes)
	n, err := file.Write(bytes)
	handleError(err)

	if n != numBytes {
		file.Close()
		return nil, errors.New("bytes written and deisred size do not match")
	}

	return file, nil
}

func NewFileFetcher(client Client, numWorkers int64, metaData *models.ResultFile) *FileFetcher {
	return &FileFetcher{client: client, numWorkers: numWorkers, shardLen: 100, meta: metaData}
}
