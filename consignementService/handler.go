package main

import (
	"context"
	"log"

	protoB "github.com/zjjt/shippingGo/consignementService/proto/consignement"
	vesselProto "github.com/zjjt/shippingGo/vesselService/proto/vessel"
)

// Service should implement all of the methods to satisfy the service
// we defined in our protobuf definition. You can check the interface
// in the generated code itself for the exact method signatures etc
// to give you a better idea.
// it is our handler
type service struct {
	repo         repository
	vesselClient vesselProto.VesselService
}

func newService(repo *MongoRepository, vesselClient vesselProto.VesselService) *service {
	return &service{repo, vesselClient}
}

//CreateConsignement API from our grpc service
func (serv *service) CreateConsignement(ctx context.Context, req *protoB.Consignement, res *protoB.Response) error {
	vesselResponse, err := serv.vesselClient.FindAvailable(context.Background(), &vesselProto.Specification{
		MaxWeight: req.Weight,
		Capacity:  int32(len(req.Containers)),
	})
	if err != nil {
		log.Println(err)
		return err
	}
	// we store in the VesselId as the vesseilId we got back from our vessel service
	req.VesselId = vesselResponse.Vessel.Id
	//save consignement in DB
	if err := serv.repo.Create(ctx, MarshalConsignement(req)); err != nil {
		log.Println(err)
		return err
	}

	res.Created = true
	res.Consignement = req
	return nil
}

//GetConsignements get all consignements
func (serv *service) GetConsignements(ctx context.Context, req *protoB.GetRequest, res *protoB.Response) error {
	consignements, err := serv.repo.GetAll(ctx)
	if err != nil {
		log.Println(err)
		return err
	}
	res.Consignements = UnmarshalConsignementCollection(consignements)
	return nil
}
