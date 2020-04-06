package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/micro/go-micro/v2"
	pb "github.com/zjjt/shippingGo/vesselService/proto/vessel"
)

//Repository - interface of what our datastore should be
type repository interface {
	FindAvailable(*pb.Specification) (*pb.Vessel, error)
}

//VesselRepository - the actual repository implementig our interface methods
type VesselRepository struct {
	Vessels []*pb.Vessel
}

//FindAvailable - checks a specification against a slice of vessels
//if the capacity and max_weight are below the vessel's then we return that vessel
func (repo *VesselRepository) FindAvailable(spec *pb.Specification) (*pb.Vessel, error) {
	for _, v := range repo.Vessels {
		if spec.Capacity <= v.Capacity && spec.MaxWeight <= v.MaxWeight {
			return v, nil
		}
	}
	return nil, errors.New("no vessel found bythat spec")
}

//NewVesselRepository - creates an empty repository
func NewVesselRepository(vessels []*pb.Vessel) *VesselRepository {
	return &VesselRepository{Vessels: vessels}
}

//our GRPC service handler
type service struct {
	repo repository
}

func newService(repo *VesselRepository) *service {
	return &service{repo}
}
func (service *service) FindAvailable(ctx context.Context, req *pb.Specification, res *pb.Response) error {
	vessel, err := service.repo.FindAvailable(req)
	if err != nil {
		return err
	}
	res.Vessel = vessel
	return nil
}

func main() {
	vessels := []*pb.Vessel{
		&pb.Vessel{Id: "vessel001", Name: "Boaty McBoatface", MaxWeight: 200000, Capacity: 500},
	}
	repo := NewVesselRepository(vessels)
	server := micro.NewService(micro.Name("shippingGo.service.vessel"))
	server.Init()
	//registering our service
	pb.RegisterVesselServiceHandler(server.Server(), newService(repo))
	//run the server
	if err := server.Run(); err != nil {
		fmt.Println(err)
	}
}
