package services

import (
	"fmt"
	ps "partitions_searcher/proto/partitions_searcher"
	"partitions_searcher/settings"
	"sync"
)

func GetRecordsInPartitions(req *ps.GetRecordsRequest) (result []*ps.Record, errorFindRecords error) {
	partitionsCount := len(req.PartitionsArray)

	taskChannel := make(chan []ps.Record, partitionsCount)
	taskCountChannel := make(chan int, partitionsCount)
	var wg sync.WaitGroup
	resultCount := 0

	for _, partition := range req.PartitionsArray {
		wg.Add(1)
		go func(partition string) {
			defer wg.Done()
			models, err := settings.GetEventByPartition(partition, req.Query)
			if err != nil {
				errorFindRecords = err
			}
			taskChannel <- models
			taskCountChannel <- len(models)
		}(partition)
	}
	wg.Wait()

	for i := 0; i < partitionsCount; i++ {
		resultCount += <-taskCountChannel
	}

	fmt.Println(len(taskChannel))
	fmt.Println(resultCount)

	result = make([]*ps.Record, resultCount)
	iterator := 0
	for models := range taskChannel {
		if len(taskChannel) == 0 {
			break
		}
		for _, model := range models {
			result[iterator] = &model
			iterator++
		}
	}
	return
}
