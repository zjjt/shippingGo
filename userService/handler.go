package main

import (
	"context"
	"errors"
	"log"

	pb "github.com/zjjt/shippingGo/userService/proto/user"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repo         repository
	tokenService Authable
}

func newUserService(repo repository, tokenService Authable) *service {
	return &service{repo, tokenService}
}

//Get - retrieves a single user
func (s *service) Get(ctx context.Context, req *pb.User, res *pb.Response) error {
	user, err := s.repo.Get(req.Id)
	if err != nil {
		return err
	}
	res.User = user
	return nil
}

//GetAll -returns a slice of users
func (s *service) GetAll(ctx context.Context, req *pb.Request, res *pb.Response) error {
	users, err := s.repo.GetAll()
	if err != nil {
		return err
	}
	res.Users = users
	return nil
}

//Auth - check against the db if the user exist or not and retrieves a token
func (s *service) Auth(ctx context.Context, req *pb.User, res *pb.Token) error {
	log.Println("login in with:", req.Email, req.Password)
	user, err := s.repo.GetByEmail(req.Email)
	log.Println(user)
	if err != nil {
		return err
	}
	// Compare the password with the hashed password stored in database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return err
	}
	token, err := s.tokenService.Encode(user)
	if err != nil {
		return err
	}
	res.Token = token
	return nil
}

//Create -creates a user in db
func (s *service) Create(ctx context.Context, req *pb.User, res *pb.Response) error {
	//generate hash version of user password
	hashPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	req.Password = string(hashPass)
	if err := s.repo.Create(req); err != nil {
		return err
	}
	res.User = req
	return nil
}

//ValidateToken -check if the token is valid
func (s *service) ValidateToken(ctx context.Context, req *pb.Token, res *pb.Token) error {
	//decode token
	claims, err := s.tokenService.Decode(req.Token)
	if err != nil {
		return err
	}
	log.Println(claims)
	if claims.User.Id == "" {
		return errors.New("invalid user")
	}
	res.Valid = true
	return nil
}
