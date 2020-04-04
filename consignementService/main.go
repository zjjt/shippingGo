package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/micro/go-micro/v2"
	protoB "github.com/zjjt/shippingGo/consignementService/proto/consignement"
)

const (
	port = ":50051"
)

type repository interface {
	Create(*protoB.Consignement) (*protoB.Consignement, error)
	GetAll() []*protoB.Consignement
}

//Repository Simulates a datastore of some kind
type Repository struct {
	mu            sync.RWMutex
	consignements []*protoB.Consignement
}

//
func newRepository() *Repository {
	return &Repository{}
}

//Create a new consignement
func (repo *Repository) Create(consignement *protoB.Consignement) (*protoB.Consignement, error) {
	repo.mu.Lock()
	updated := append(repo.consignements, consignement)
	repo.consignements = updated
	repo.mu.Unlock()
	return consignement, nil
}

//GetAll gets all the consignement
func (repo *Repository) GetAll() []*protoB.Consignement {
	return repo.consignements
}

// Service should implement all of the methods to satisfy the service
// we defined in our protobuf definition. You can check the interface
// in the generated code itself for the exact method signatures etc
// to give you a better idea.
type service struct {
	repo repository
}

func newService(repo *Repository) *service {
	return &service{repo}
}

//CreateConsignement API from our grpc service
func (serv *service) CreateConsignement(ctx context.Context, req *protoB.Consignement, res *protoB.Response) error {
	//save consignement in DB
	consignement, err := serv.repo.Create(req)
	if err != nil {
		return err
	}

	res.Consignement = consignement
	return nil
}

//GetConsignements get all consignements
func (serv *service) GetConsignements(ctx context.Context, req *protoB.GetRequest, res *protoB.Response) error {
	consignements := serv.repo.GetAll()
	res.Consignements = consignements
	return nil
}

func main() {
	repo := newRepository()

	//Create a grpc server with go micro
	//the name must match the package name given in the protobuf file for service discovery
	server := micro.NewService(micro.Name("shippingGo.service.consignement"))
	//will parse the command line flags
	server.Init()
	//Registering our service hangler
	protoB.RegisterShippingServiceHandler(server.Server(), newService(repo))
	if err := server.Run(); err != nil {
		fmt.Println(err)
	}

}
