package main

import (
	"context"
	"fmt"
	pb "github.com/krishnareddyakkala/golang_practice/proto_buf/deos_ingestion_v1alpa1_stream"
	"io"
	"log"
	"math/rand"
	"net"
	"time"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedIngestionServer
	ingestionList []*pb.IngestionRequest
}

func (s *server) PostSingleIngestion(ctx context.Context, in *pb.IngestionRequest) (*pb.IngestionResponse, error) {
	log.Printf("Received: %v", in.GetAccountId())
	res := &pb.IngestionResponse{}
	res.Response = fmt.Sprintf("Received account ID: %s", in.GetAccountId())
	return res, nil
}

func (s *server) PostStreamIngestion(stream pb.Ingestion_PostStreamIngestionServer) error {
	var count int32 = 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			if err := stream.SendAndClose(&pb.IngestionCount{IngestionCount: count}); err != nil {
				return err
			}
			return nil
		}
		count += 1
		log.Printf("Received: %v", req.GetAccountId())
		s.ingestionList = append(s.ingestionList, req)
	}
}

func (s *server) ListIngestionStream(req *pb.IngestionRequest, stream pb.Ingestion_ListIngestionStreamServer) error {
	fmt.Println("ingestion ListIngestionStream : ", req.GetAccountId())
	for _, ingestion := range s.ingestionList {
		if ingestion.GetAccountId() != req.GetAccountId() {
			continue
		}
		r := &pb.IngestionResponse{}
		r.AccountId = ingestion.GetAccountId()
		r.Response = fmt.Sprintf("Received account ID: %s", ingestion.GetAccountId())
		if err := stream.Send(r); err != nil {
			return err
		}
	}

	return nil
}

func main() {

	log.Println("Starting Server")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterIngestionServer(s, &server{})
	log.Println("Server started successfully.")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
