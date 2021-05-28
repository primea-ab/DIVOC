package fetcher

import (
	"fmt"
	"log"
	"math"
	"sync"

	"divoc.primea.se/app/client/shardclient"
	"divoc.primea.se/models"
	"divoc.primea.se/util"
)

type FileFetcher struct {
	client          shardclient.Client
	numWorkers      int64
	shardLen        int64
	meta            *models.ResultFile
	progressChannel *chan float64
}

func handleError(err error) {
	if err != nil {
		log.Fatalf("Something went wrong: %v\n", err)
	}
}

func (a *FileFetcher) Download() error {
	defer close(*a.progressChannel)
	filePath := fmt.Sprintf("./%s/%s", util.RootShareFolder, a.meta.Names[0])
	file, err := a.getFile(int(a.meta.Size), filePath)
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

	go a.writeWorker(dataChan, file, &wg, a.shardLen, float64(numShards))

	for i = 0; i < a.numWorkers; i += 1 {
		go a.fetchWorker(dataChan, shardChan, i)
	}

	wg.Wait()

	return nil
}

func New(client shardclient.Client, numWorkers int64, metaData *models.ResultFile, progressChannel *(chan float64)) *FileFetcher {
	return &FileFetcher{client: client, numWorkers: numWorkers, shardLen: 100, meta: metaData, progressChannel: progressChannel}
}
