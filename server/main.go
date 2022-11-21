package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	ps "partitions_searcher/proto/partitions_searcher"
)

type GetRecordsServiceServer struct {
	ps.GetRecordsServiceServer
}

func (s *GetRecordsServiceServer) Get(ctx context.Context,
	req *ps.GetRecordsRequest) (*ps.GetRecordsResponse, error) {

	var err error
	response := new(ps.GetRecordsResponse)

	response.Record = make([]*ps.Record, 10)

	return response, err
}

func main() {
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
