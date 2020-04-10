package main

import (
	"github.com/dgrijalva/jwt-go"
	pb "github.com/zjjt/shippingGo/userService/proto/user"
)

var (
	key = []byte("it's}a{fucking*Secret%lol")
)

// CustomClaims is our custom metadata, which will be hashed
// and sent as the second segment in our JWT
type CustomClaims struct {
	User *pb.User
	jwt.StandardClaims
}

type Authable interface {
	Decode(token string) (*CustomClaims, error)
	Encode(user *pb.User) (string, error)
}

type TokenService struct {
	repo repository
}

func newtokenService(repo repository) *TokenService {
	return &TokenService{repo}
}
func (srv *TokenService) Decode(token string) (*CustomClaims, error) {
	//parse the token
	tokentype, err := jwt.ParseWithClaims(string(key), &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	//validate the token and return the customclaims
	if claims, ok := tokentype.Claims.(*CustomClaims); ok && tokentype.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func (srv *TokenService) Encode(user *pb.User) (string, error) {
	claims := CustomClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    "shippingGo.service.user",
		},
	}
	//create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(key)
}
