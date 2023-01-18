package services

import (
	ps "partitions_searcher/proto/partitions_searcher"
	"sort"
	"testing"
)

func TestSortRecords(t *testing.T) {
	var testCoreService CoreService
	testRecords := make([]*ps.Record, 10)
	testIds := []int{1, 5, 4, 4, 2, 2, 4, 66, 6345, 2}

	for i := 0; i < 10; i++ {
		testRecords[i] = new(ps.Record)
		testRecords[i].Id = int64(testIds[i])
	}

	testCoreService.SortDirection = true
	result := testCoreService.mergeSort(testRecords)

	sort.Ints(testIds)

	for i := 0; i < 10; i++ {
		if result[i].Id != int64(testIds[i]) {
			t.Fatalf("%v != %v", result[i].Id, testIds[i])
		}
	}
}
