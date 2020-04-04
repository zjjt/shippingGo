package main

import (
	"context"
	"log"
	"net"
	"sync"

	protoB "github.com/zjjt/shippingGo/consignementService/proto/consignement"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

//CreateConsignement APIfrom our grpc service
func (serv *service) CreateConsignement(ctx context.Context, req *protoB.Consignement) (*protoB.Response, error) {
	//save consignement in DB
	consignement, err := serv.repo.Create(req)
	if err != nil {
		return nil, err
	}
	return &protoB.Response{Created: true, Consignement: consignement}, nil
}

//GetConsignements get all consignements
func (serv *service) GetConsignements(ctx context.Context, req *protoB.GetRequest) (*protoB.Response, error) {
	consignements := serv.repo.GetAll()
	return &protoB.Response{Consignements: consignements}, nil
}

func main() {
	repo := newRepository()
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen to port %s \n err %v", port, err)
	}
	//Create a grpc server
	server := grpc.NewServer()
	//Registering our service
	protoB.RegisterShippingServiceServer(server, newService(repo))
	reflection.Register(server)
	log.Println("Running on port:", port)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
