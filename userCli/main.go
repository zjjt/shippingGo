package main

import (
	"context"
	"log"
	"os"

	"github.com/micro/go-micro/v2"
	microclient "github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/config/cmd"
	pb "github.com/zjjt/shippingGo/userService/proto/user"
)

func main() {

	cmd.Init()
	client := pb.NewUserService("shippingGo.service.user", microclient.DefaultClient)
	//utilisation de la ligne de commande pour passer des parametres via l'outil de go-micro
	service := micro.NewService(micro.Name("shippingGo.user.cli"))
	//start the service
	service.Init()
	name := "thibaut"
	email := "zehijean1988@gmail.com"
	password := "0123456789"
	company := "TECHNIKING SA"
	//here we call our user service
	r, err := client.Create(context.TODO(), &pb.User{
		Name:     name,
		Email:    email,
		Password: password,
		Company:  company,
	})
	if err != nil {
		log.Fatalf("couldnt create user %v", err)
	}
	log.Printf("Created: %s", r.User.Id)
	getAll, err := client.GetAll(context.Background(), &pb.Request{})
	if err != nil {
		log.Fatalf("Couldnt retrieve list of users %v", err)
	}
	for _, v := range getAll.Users {
		log.Println(v)
	}
	authResponse, err := client.Auth(context.Background(), &pb.User{
		Email:    email,
		Password: password,
	})

	if err != nil {
		log.Fatalf("Could not authenticate user: %s error: %v\n", email, err)
	}
	log.Printf("Your access token is: %s \n", authResponse.Token)

	os.Exit(0)

}
