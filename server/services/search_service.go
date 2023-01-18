package services

import (
	ps "partitions_searcher/proto/partitions_searcher"
	"partitions_searcher/server/custom"
	"sync"
)

type CoreService struct {
	Error         error
	SortDirection bool
}

func (s *CoreService) GetRecordsInPartitions(req *ps.GetRecordsRequest) (result []*ps.Record) {
	partitionsCount := len(req.PartitionsArray)

	modelChannel := make(chan []ps.Record, partitionsCount)
	modelCountChannel := make(chan int, partitionsCount)
	resultCount := 0
	var wg sync.WaitGroup
	for _, partition := range req.PartitionsArray {
		wg.Add(1)
		go func(partition string) {
			defer wg.Done()
			models, err := custom.GetEventByPartition(partition, req.Query)
			if err != nil {
				s.Error = err
			}
			modelChannel <- models
			modelCountChannel <- len(models)
		}(partition)
	}
	wg.Wait()

	for i := 0; i < partitionsCount; i++ {
		resultCount += <-modelCountChannel
	}
	close(modelCountChannel)

	result = make([]*ps.Record, resultCount)
	iterator := 0
	for i := 0; i < partitionsCount; i++ {
		if len(modelChannel) == 0 {
			break
		}
		models := <-modelChannel
		for j := 0; j < len(models); j++ {
			result[iterator] = &models[j]
			iterator++
		}
	}
	close(modelChannel)

	return s.mergeSort(result)
}

func (s *CoreService) mergeSort(r []*ps.Record) []*ps.Record {
	if len(r) <= 1 {
		return r
	}
	var wg sync.WaitGroup
	n := len(r) / 2
	wg.Add(2)

	var r1 []*ps.Record
	var r2 []*ps.Record

	go func() {
		r1 = s.mergeSort(r[:n])
		wg.Done()
	}()

	go func() {
		r2 = s.mergeSort(r[n:])
		wg.Done()
	}()

	wg.Wait()
	return merge(r1, r2, s.SortDirection)
}

func merge(l, r []*ps.Record, sortDirection bool) []*ps.Record {
	ret := make([]*ps.Record, 0, len(l)+len(r))
	for len(l) > 0 || len(r) > 0 {
		if len(l) == 0 {
			return append(ret, r...)
		}
		if len(r) == 0 {
			return append(ret, l...)
		}
		if custom.GetDifference(l[0], r[0], sortDirection) {
			ret = append(ret, l[0])
			l = l[1:]
		} else {
			ret = append(ret, r[0])
			r = r[1:]
		}
	}
	return ret
}
