package fetcher

import (
	"fmt"
	"os"
	"sync"

	"divoc.primea.se/util"
)

type partialResult struct {
	payload []byte
	id      int64
}

func (a *FileFetcher) fetchWorker(dataChan chan<- partialResult, shardIdChan <-chan int64, workerId int64) {
	for id := range shardIdChan {
		res, err := a.client.GetShard(id, *a.meta)
		handleError(err)
		dataChan <- partialResult{id: id, payload: res}
	}
}

func (a *FileFetcher) writeWorker(dataChan <-chan partialResult, file *os.File, wg *sync.WaitGroup, numShards float64) {
	var numDownloads float64 = 0

	for data := range dataChan {
		_, err := file.WriteAt(data.payload, data.id*util.ChunkByteSize)
		if err != nil {
			fmt.Printf("Failed to write dat to file: %+v\n", err)
		}
		numDownloads += 1
		a.Progress = numDownloads / numShards
		wg.Done()
	}
}
