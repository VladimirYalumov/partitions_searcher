package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	ps "partitions_searcher/proto/partitions_searcher"
)

func main() {
	fmt.Println(getTasks())
}

func getTasks() ([]*ps.Record, error) {
	conn, _ := grpc.Dial("127.0.0.1:8088", grpc.WithInsecure())
	partitions := []string{"tasks_20"}
	client := ps.NewGetRecordsServiceClient(conn)

	resp, err := client.Get(
		context.Background(),
		&ps.GetRecordsRequest{
			PartitionsArray: partitions,
			SortField:       "id",
			SortDirection:   true,
			Query:           " where ;",
		},
	)

	if err != nil {
		return []*ps.Record{}, err
	}

	return resp.GetRecord(), nil
}
