package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/micro/go-micro/v2"
	protoB "github.com/zjjt/shippingGo/consignementService/proto/consignement"
	vesselProto "github.com/zjjt/shippingGo/vesselService/proto/vessel"
)

const (
	defaultDbHost = "datastore:27017"
)

func main() {

	//Create a grpc server with go micro
	//the name must match the package name given in the protobuf file for service discovery
	server := micro.NewService(micro.Name("shippingGo.service.consignement"))
	//will parse the command line flags
	server.Init()
	uri := os.Getenv("DB_HOST")
	log.Println("uri is ", uri)
	if uri == "" {
		uri = defaultDbHost
	}
	//connect to the database
	client, err := CreateDBClient(context.Background(), uri, 0)
	if err != nil {
		log.Panic(err)
	}
	defer client.Disconnect(context.Background())
	//create a collection in the database
	consignementCollection := client.Database("shippinggo").Collection("consignements")
	repository := NewRepository(consignementCollection)
	vesselClient := vesselProto.NewVesselService("shippingGo.service.vessel", server.Client())
	//Registering our service hangler
	handler := newService(repository, vesselClient)
	protoB.RegisterShippingServiceHandler(server.Server(), handler)
	if err := server.Run(); err != nil {
		fmt.Println(err)
	}

}
