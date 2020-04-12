package main

import (
	"time"

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

//Decode -check the token and see if it is valid and retrieves the details encoded within
func (srv *TokenService) Decode(token string) (*CustomClaims, error) {
	//parse the token
	tokentype, err := jwt.ParseWithClaims(string(token), &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	//validate the token and return the customclaims
	if claims, ok := tokentype.Claims.(*CustomClaims); ok && tokentype.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

//Encode -encode the data and creates the jwt
func (srv *TokenService) Encode(user *pb.User) (string, error) {
	exprireTime := time.Now().Add(time.Hour * 72).Unix()
	claims := CustomClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: exprireTime,
			Issuer:    "shippingGo.service.user",
		},
	}
	//create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(key)
}
