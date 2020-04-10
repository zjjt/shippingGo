package main

import (
	"fmt"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/util/log"
	pb "github.com/zjjt/shippingGo/userService/proto/user"
)

func main() {
	db, err := CreateDBConnection()
	defer db.Close()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
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
	pb.RegisterUserServiceHandler(service.Server(), newUserService(repo, tokenservice))
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
