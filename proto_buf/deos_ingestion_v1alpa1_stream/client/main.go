// Package main implements a client for Greeter service.
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

	pb "github.com/krishnareddyakkala/golang_practice/proto_buf/deos_ingestion_v1alpa1_stream"
	"google.golang.org/grpc"
)

const (
	address         = "localhost:50051"
	ingestionLength = 20
)

func ingestStream(client pb.IngestionClient, ingestionList []pb.IngestionRequest) {
	log.Println("testing PostStreamIngestion ")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.PostStreamIngestion(ctx)
	if err != nil {
		log.Fatalf("%v.RecordRoute(_) = _, %v", client, err)
	}
	for _, ingestion := range ingestionList {
		log.Println("sending AccountID: ", ingestion.GetAccountId())
		if err := stream.Send(&ingestion); err != nil {
			log.Fatalf("%v.Send(%s) = %v", stream, ingestion.GetAccountId(), err)
		}
	}
	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
	}
	log.Println("summary: ", reply.IngestionCount)
}

func ingestSingleItem(client pb.IngestionClient, ingestion pb.IngestionRequest) {
	log.Println("testing PostSingleIngestion with AccountID: ", ingestion.GetAccountId())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r, err := client.PostSingleIngestion(ctx, &ingestion)
	if err != nil {
		log.Fatalf("PostSingleIngestion failed: %v", err)
	}
	log.Println("PostSingleIngestion response: ", r.GetResponse())
}

func listIngestedStream(client pb.IngestionClient, ingestionReq pb.IngestionRequest) (res []*pb.IngestionResponse) {

	log.Println("testing ListIngestionStream")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.ListIngestionStream(ctx, &ingestionReq)
	if err != nil {
		log.Fatalf("%v.ListIngestionStream(_) = _, %v", client, err)
	}
	for {
		ingestionRes, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.listIngestedStream(_) = _, %v", client, err)
		}
		log.Println(ingestionRes)
		res = append(res, ingestionRes)
	}

	return res
}

func createIngestionList() []pb.IngestionRequest {
	ingestionList := make([]pb.IngestionRequest, ingestionLength)
	for count := 0; count < ingestionLength; count++ {
		ingestionList[count].AccountId = fmt.Sprintf("AccountID-%v", count)
		ingestionList[count].CompressionType = pb.CompressionType(count)
	}
	return ingestionList
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewIngestionClient(conn)

	ingestionList := createIngestionList()

	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()
	//req := pb.IngestionRequest{}
	//req.AccountId = "AccountID-1"
	//req.CompressionType = pb.CompressionType(32)
	//r, err := c.Post(ctx, &req)
	//if err != nil {
	//	log.Fatalf("could not greet: %v", err)
	//}
	//log.Printf("Response: %s", r.GetResponse())

	ingestSingleItem(c, ingestionList[0])
	ingestStream(c, ingestionList)
	listIngestedStream(c, ingestionList[rand.Intn(ingestionLength-1)])

}
