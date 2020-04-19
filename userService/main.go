package main

import (
	"fmt"
	"log"

	"github.com/micro/go-micro/v2"
	pb "github.com/zjjt/shippingGo/userService/proto/user"
)

func main() {
	db, err := CreateDBConnection()
	defer db.Close()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	} else {
		log.Println("Connected to DB successfully")
	}
	// Automatically migrates the user struct
	// into database columns/types etc. This will
	// check for changes and migrate them each time
	// this service is restarted.
	db.AutoMigrate(&pb.User{})
	repo := newUserRepository(db)
	tokenservice := newtokenService(repo)
	service := micro.NewService(micro.Name("shippingGo.service.user"))
	service.Init()
	//get an instance of the event broker
	pubsub := service.Server().Options().Broker
	if err := pubsub.Init(); err != nil {
		log.Fatalf("Broker Init error: %v", err)
	}
	if err := pubsub.Connect(); err != nil {
		log.Fatalf("Broker Connect error: %v", err)
	}
	pb.RegisterUserServiceHandler(service.Server(), newUserService(repo, tokenservice, pubsub))
	if err := service.Run(); err != nil {
		theerror := fmt.Sprintf("%v --from UserService", err)
		fmt.Println(theerror)
	}
}
