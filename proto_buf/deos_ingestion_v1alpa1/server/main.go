package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/krishnareddyakkala/golang_practice/proto_buf/deos_ingestion_v1alpa1"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedIngestionServer
}

func (s *server) Post(ctx context.Context, in *pb.IngestionRequest) (*pb.IngestionResponse, error) {
	log.Printf("Received: %v", in.GetAccountId())
	res := &pb.IngestionResponse{}
	res.Response = fmt.Sprintf("Received account ID: %s", in.GetAccountId())
	return res, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterIngestionServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
