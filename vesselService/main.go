package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/micro/go-micro/v2"
	pb "github.com/zjjt/shippingGo/vesselService/proto/vessel"
)

const (
	defaultDbHost = "datastore:27017"
)

func createDummyVessels(repo repository) {
	vessels := []*Vessel{
		{ID: "vessel001", Name: "Boaty McBoatface", MaxWeight: 200000, Capacity: 500},
	}
	for _, v := range vessels {
		log.Printf("v is %v\n", v)
		err := repo.Create(context.Background(), v)
		if err != nil {
			log.Println("erreur survenue lors de la creation dun vessel ", err)
		}
	}
}
func main() {

	server := micro.NewService(micro.Name("shippingGo.service.vessel"))
	server.Init()
	uri := os.Getenv("DB_HOST")
	if uri == "" {
		uri = defaultDbHost
	}
	log.Println("uri is ", uri)
	//connect to the database
	client, err := CreateDBClient(context.Background(), uri, 0)
	if err != nil {
		theerror := fmt.Sprintf("%v --from VesselService", err)
		log.Panic(theerror)
	}
	defer client.Disconnect(context.Background())
	//create a collection in the database
	vesselCollection := client.Database("shippinggo").Collection("vessels")
	repository := NewVesselRepository(vesselCollection)
	handler := newService(repository)
	createDummyVessels(repository)
	//registering our service
	pb.RegisterVesselServiceHandler(server.Server(), handler)
	//run the server
	if err := server.Run(); err != nil {
		theerror := fmt.Sprintf("%v --from ConsignementService", err)
		fmt.Println(theerror)
	}
}
