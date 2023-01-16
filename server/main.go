package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	ps "partitions_searcher/proto/partitions_searcher"
	"partitions_searcher/server/services"
	"partitions_searcher/settings"
)

type GetRecordsServiceServer struct {
	ps.GetRecordsServiceServer
}

func (s *GetRecordsServiceServer) Get(ctx context.Context,
	req *ps.GetRecordsRequest) (*ps.GetRecordsResponse, error) {

	var err error
	response := new(ps.GetRecordsResponse)

	response.Records, err = services.GetRecordsInPartitions(req)

	return response, err
}

func main() {
	settings.Init("config.yml")
	server := grpc.NewServer()

	instance := new(GetRecordsServiceServer)

	ps.RegisterGetRecordsServiceServer(server, instance)

	listener, err := net.Listen("tcp", ":8088")
	if err != nil {
		log.Fatal("Unable to create grpc listener:", err)
	}
	log.Println("Listen port 8088")
	if err = server.Serve(listener); err != nil {
		log.Fatal("Unable to start server:", err)
	}
}
