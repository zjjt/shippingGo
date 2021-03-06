package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/micro/go-micro/v2/broker"
	pb "github.com/zjjt/shippingGo/userService/proto/user"
	"golang.org/x/crypto/bcrypt"
)

//event to be sent
const topic = "user.created"

type service struct {
	repo         repository
	tokenService Authable
	PubSub       broker.Broker
}

func newUserService(repo repository, tokenService Authable, pubsub broker.Broker) *service {
	return &service{repo, tokenService, pubsub}
}

//Get - retrieves a single user
func (s *service) Get(ctx context.Context, req *pb.User, res *pb.Response) error {
	user, err := s.repo.Get(req.Id)
	if err != nil {
		theerror := fmt.Sprintf("%v --from UserService", err)
		return errors.New(theerror)
	}
	res.User = user
	return nil
}

//GetAll -returns a slice of users
func (s *service) GetAll(ctx context.Context, req *pb.Request, res *pb.Response) error {
	users, err := s.repo.GetAll()
	if err != nil {
		theerror := fmt.Sprintf("%v --from UserService", err)
		return errors.New(theerror)
	}
	res.Users = users
	return nil
}

//Auth - check against the db if the user exist or not and retrieves a token
func (s *service) Auth(ctx context.Context, req *pb.User, res *pb.Token) error {
	log.Println("login in with:", req.Email, req.Password)
	user, err := s.repo.GetByEmail(req.Email)
	log.Println("user is ", user)
	if err != nil {
		theerror := fmt.Sprintf("%v --from UserService", err)
		return errors.New(theerror)
	}
	// Compare the password with the hashed password stored in database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		theerror := fmt.Sprintf("%v --from UserService", err)
		return errors.New(theerror)
	}
	token, err := s.tokenService.Encode(user)
	if err != nil {
		theerror := fmt.Sprintf("%v --from UserService", err)
		return errors.New(theerror)
	}
	res.Token = token
	return nil
}

//Create -creates a user in db
func (s *service) Create(ctx context.Context, req *pb.User, res *pb.Response) error {
	//generate hash version of user password
	hashPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		theerror := fmt.Sprintf("%v --from UserService", err)
		return errors.New(theerror)
	}
	req.Password = string(hashPass)
	if err := s.repo.Create(req); err != nil {
		theerror := fmt.Sprintf("%v --from UserService", err)
		return errors.New(theerror)
	}
	res.User = req
	if err := s.publishEvent(req); err != nil {
		return err
	}
	return nil
}

//ValidateToken -check if the token is valid
func (s *service) ValidateToken(ctx context.Context, req *pb.Token, res *pb.Token) error {
	//decode token
	claims, err := s.tokenService.Decode(req.Token)
	if err != nil {
		theerror := fmt.Sprintf("%v --from UserService", err)
		return errors.New(theerror)
	}
	//log.Println(claims)
	if claims.User.Id == "" {
		return errors.New("invalid user --from UserService")
	}
	res.Valid = true
	return nil
}
func (s *service) publishEvent(user *pb.User) error {
	//when sending an event we have to serialize it to bytes
	//we are sending to our ecosystem the event user.created with the details
	//concerning that user
	body, err := json.Marshal(user)
	if err != nil {
		theerror := fmt.Sprintf("%v --from UserService", err)
		return errors.New(theerror)
	}

	//create a broker message
	msg := &broker.Message{
		Header: map[string]string{
			"id": user.Id,
		},
		Body: body,
	}
	//publish the message to the broker
	if err := s.PubSub.Publish(topic, msg); err != nil {
		theerror := fmt.Sprintf("%v --from UserService", err)
		log.Printf("[PUB] failed %s\n", theerror)
	}
	return nil
}
