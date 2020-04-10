package main

import (
	"context"
	"log"

	pb "github.com/zjjt/shippingGo/vesselService/proto/vessel"
)

//our GRPC service handler
type service struct {
	repo repository
}

func newService(repo *VesselRepository) *service {
	return &service{repo}
}
func (service *service) FindAvailable(ctx context.Context, req *pb.Specification, res *pb.Response) error {
	vessel, err := service.repo.FindAvailable(ctx, MarshalSpecification(req))
	if err != nil {
		return err
	}
	res.Vessel = UnmarshalVessel(vessel)
	return nil
}
func (service *service) Create(ctx context.Context, req *pb.Vessel, res *pb.Response) error {
	err := service.repo.Create(ctx, MarshalVessel(req))
	if err != nil {
		log.Printf("Couldn't create a new vessel %v \n", req)
	}
	return err
}
