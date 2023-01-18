package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"partitions_searcher/example/models"
	ps "partitions_searcher/proto/partitions_searcher"
	"strconv"
)

func main() {
	tasks, _ := getTasks()
	id := ""
	for _, task := range tasks {
		id = strconv.Itoa(int(task.Id))
		fmt.Println("ID: " + id + " Title: " + task.Title + " Description: " + task.Description)
	}
}

func getTasks() (task []models.Task, err error) {
	conn, _ := grpc.Dial("127.0.0.1:8088", grpc.WithInsecure())
	partitions := []string{
		"tasks_2022_01",
		"tasks_2022_02",
		"tasks_2022_03",
		"tasks_2022_04",
		"tasks_2022_05",
		"tasks_2022_06",
		"tasks_2022_07",
		"tasks_2022_08",
		"tasks_2022_09",
		"tasks_2022_10",
		"tasks_2022_11",
		"tasks_2022_12",
	}
	client := ps.NewGetRecordsServiceClient(conn)

	resp, err := client.Get(
		context.Background(),
		&ps.GetRecordsRequest{
			PartitionsArray: partitions,
			SortDirection:   true,
			Query:           "where description = 'test description' order by id;",
		},
	)

	if err != nil {
		return
	}

	tasks := make([]models.Task, len(resp.GetRecords()))

	for i, record := range resp.GetRecords() {
		tasks[i].Id = record.GetId()
		tasks[i].Title = record.GetTitle()
		tasks[i].Description = record.GetDescription()
	}

	return tasks, nil
}
