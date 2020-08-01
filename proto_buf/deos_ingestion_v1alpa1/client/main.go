// Package main implements a client for Greeter service.
package main

import (
	"context"
	"log"
	"time"

	pb "github.com/krishnareddyakkala/golang_practice/proto_buf/deos_ingestion_v1alpa1"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewIngestionClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := pb.IngestionRequest{}
	req.AccountId = "AccountID-1"
	req.CompressionType = pb.CompressionType(32)

	r, err := c.Post(ctx, &req)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Response: %s", r.GetResponse())

}
