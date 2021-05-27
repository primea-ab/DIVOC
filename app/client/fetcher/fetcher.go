package fetcher

import (
	"errors"
	"log"
	"math"
	"os"
	"sync"
)

type FileFetcher struct {
	client     Client
	filename   string
	hash       string
	shardLen   int64
	numWorkers int64
}

func handleError(err error) {
	if err != nil {
		log.Fatalf("Something went wrong: %v\n", err)
	}
}

func (a *FileFetcher) Download() error {
	data, err := a.client.GetFileMetaData(a.filename)
	handleError(err)

	file, err := a.getFile(data.Size)
	handleError(err)
	defer file.Close()

	numShards := int64(math.Ceil(float64(data.Size) / float64(a.shardLen)))

	shardChan := make(chan int64, numShards)
	dataChan := make(chan partialResult, a.numWorkers)
	var wg sync.WaitGroup
	wg.Add(int(numShards))

	var i int64
	for i = 0; i < numShards; i += 1 {
		shardChan <- i
	}

	go a.startWriteWorker(dataChan, file, &wg)

	for i = 0; i < a.numWorkers; i += 1 {
		go a.startFetchWorker(dataChan, shardChan)
	}

	wg.Wait()

	return nil
}

func (a *FileFetcher) getFile(numBytes int) (*os.File, error) {
	file, err := os.Open(a.filename)
	if err != nil {
		return a.createEmpty(numBytes)
	}

	info, _ := file.Stat()
	if info.Size() == int64(numBytes) {
		return file, nil
	}

	return a.createEmpty(numBytes)
}

func (a *FileFetcher) createEmpty(numBytes int) (*os.File, error) {
	file, err := os.Create(a.filename)
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

func NewFileFetcher(client Client, filename string) *FileFetcher {
	return &FileFetcher{client: client, filename: filename}
}
