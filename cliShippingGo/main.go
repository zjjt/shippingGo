package main

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/micro/go-micro/v2"
	microclient "github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/config/cmd"
	"github.com/micro/go-micro/v2/metadata"
	pb "github.com/zjjt/shippingGo/consignementService/proto/consignement"
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
	cmd.Init()
	//create a client to exchange with the grpc server
	client := pb.NewShippingService("shippingGo.service.consignement", microclient.DefaultClient)
	service := micro.NewService(micro.Name("shippingGo.consignement.cli"))
	//start the service
	service.Init()
	file := defaultFilename
	var token string
	log.Println(os.Args)
	if len(os.Args) < 3 {
		log.Fatal(errors.New("Not enough arguments,expecting file and token"))
	}
	file = os.Args[1]
	token = os.Args[2]
	consignement, err := parseFile(file)
	if err != nil {
		log.Fatalf("Could not parse file: %v", err)
	}
	//Create a new context which contains our given token
	//this same context will be passed into both the calls we make
	//to our consignement service
	ctx := metadata.NewContext(context.Background(), map[string]string{"token": token})
	log.Println("content of context is: ", ctx)
	//and here we pass it into a call to create a consignement
	r, err := client.CreateConsignement(ctx, consignement)
	if err != nil {
		log.Fatalf("Could not create consignement: %v", err)
	}
	log.Printf("Created: %t", r.Created)
	//second call
	getAll, err := client.GetConsignements(ctx, &pb.GetRequest{})
	if err != nil {
		log.Fatalf("could not get all the consignements %v", err)
	}
	for _, v := range getAll.Consignements {
		log.Println(v)
	}
}
