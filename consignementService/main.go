package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/server"
	protoB "github.com/zjjt/shippingGo/consignementService/proto/consignement"
	userproto "github.com/zjjt/shippingGo/userService/proto/user"
	vesselProto "github.com/zjjt/shippingGo/vesselService/proto/vessel"
)

const (
	defaultDbHost = "datastore:27017"
)

//AuthWrapper is a higher order function that takes a HandlerFunc
//and returns a function,which takes a context,request and response interface.
//The token is extracted from the context set in our consignement-cli,
//that token is then sent over to the user service to be validated.
//if valid the call is passed along to the handler, if not an error is returned
func AuthWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, res interface{}) error {
		meta, ok := metadata.FromContext(ctx)
		if !ok {
			return errors.New("no auth meta-data found in request --from ConsignementService")
		}
		//it changed from token to Token
		token := meta["Token"]
		log.Println("Authenticated with token: ", token)
		//Checking if the token is a valid one
		authClient := userproto.NewUserService("shippingGo.service.user", client.DefaultClient)
		_, err := authClient.ValidateToken(context.Background(), &userproto.Token{
			Token: token,
		})
		if err != nil {
			theerror := fmt.Sprintf("%v --from UserService", err)
			return errors.New(theerror)
		}
		err = fn(ctx, req, res)
		return err
	}
}
func main() {

	//Create a grpc server with go micro
	//the name must match the package name given in the protobuf file for service discovery
	//we add a middleware for authentication via jwt
	server := micro.NewService(micro.Name("shippingGo.service.consignement"), micro.WrapHandler(AuthWrapper))
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
		theerror := fmt.Sprintf("%v --from ConsignementService", err)
		log.Panic(theerror)
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
		theerror := fmt.Sprintf("%v --from ConsignementService", err)
		fmt.Println(theerror)
	}

}
