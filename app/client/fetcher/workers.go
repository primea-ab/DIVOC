package fetcher

import (
	"fmt"
	"os"
	"sync"
)

type partialResult struct {
	payload []byte
	id      int64
}

func (a *FileFetcher) startFetchWorker(dataChan chan<- partialResult, shardIdChan <-chan int64, workerId int64) {
	for id := range shardIdChan {
		fmt.Printf("Fetch shard %v, from worker: %v\n", id, workerId)
		res, err := a.client.GetShard(id, *a.meta)
		handleError(err)
		dataChan <- partialResult{id: id, payload: res}
	}
}

func (a *FileFetcher) startWriteWorker(dataChan <-chan partialResult, file *os.File, wg *sync.WaitGroup, shardLen int64) {
	for data := range dataChan {
		fmt.Println("Write data to file")
		_, err := file.WriteAt(data.payload, data.id*shardLen)
		if err != nil {
			fmt.Printf("Failed to write dat to file: %+v\n", err)
		}
		wg.Done()
	}
}
