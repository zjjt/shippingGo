package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	pb "github.com/zjjt/shippingGo/proto/consignement"
	"google.golang.org/grpc"
)

const (
	address         = "localhost:50051"
	defaultFilename = "consignement.json"
)

func parseFile(file string) (*pb.Consignement, error) {
	var consignement *pb.Consignement
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &consignement)
	return consignement, err
}

func main() {
	//set up a connection to the grpc server
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("couldnt connect to the server err: %v", err)
	}
	defer conn.Close()
	//create a client to exchange with the grpc server
	client := pb.NewShippingServiceClient(conn)
	file := defaultFilename
	if len(os.Args) > 1 {
		file = os.Args[1]
	}
	consignement, err := parseFile(file)
	if err != nil {
		log.Fatalf("Could not parse file: %v", err)
	}
	r, err := client.CreateConsignement(context.Background(), consignement)
	if err != nil {
		log.Fatalf("Could not create consignement: %v", err)
	}
	log.Printf("Created: %t", r.Created)
	getAll, err := client.GetConsignements(context.Background(), &pb.GetRequest{})
	if err != nil {
		log.Fatalf("could not get all the consignements %v", err)
	}
	for _, v := range getAll.Consignements {
		log.Println(v)
	}
}
