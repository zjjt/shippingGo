package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/broker"
	pb "github.com/zjjt/shippingGo/userService/proto/user"
)

const topic = "user.created"

func main() {
	srv := micro.NewService(micro.Name("shippingGo.service.email"))
	srv.Init()
	//get the broker instance
	pubsub := srv.Server().Options().Broker
	if err := pubsub.Connect(); err != nil {
		log.Fatal(err)
	}
	//Subscribe to messages on the broker
	_, err := pubsub.Subscribe(topic, func(p broker.Event) error {
		var user *pb.User
		if err := json.Unmarshal(p.Message().Body, &user); err != nil {
			theerror := fmt.Sprintf("%v --from emailService", err)
			return errors.New(theerror)
		}
		log.Println(user)
		go sendEmail(user)
		return nil
	})
	if err != nil {
		log.Println(err)
	}
	//run the server
	if err := srv.Run(); err != nil {
		log.Println(err)
	}
}

func sendEmail(user *pb.User) error {
	log.Println("sendig email to:", user.Name)
	return nil
}
