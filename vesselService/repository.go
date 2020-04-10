package main

import (
	"context"
	"log"

	pb "github.com/zjjt/shippingGo/vesselService/proto/vessel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//Here we create our models
type Vessel struct {
	ID        string `json:"id"`
	Capacity  int32  `json:"capacity"`
	MaxWeight int32  `json:"max_weight"`
	Name      string `json:"name"`
	Available bool   `json:"available"`
	OwnerID   string `json:"owner_id"`
}
type Specification struct {
	Capacity  int32
	MaxWeight int32
}

func MarshalVessel(v *pb.Vessel) *Vessel {
	return &Vessel{
		ID:        v.Id,
		Capacity:  v.Capacity,
		MaxWeight: v.MaxWeight,
		Name:      v.Name,
		Available: v.Available,
		OwnerID:   v.OwnerId,
	}
}
func UnmarshalVessel(v *Vessel) *pb.Vessel {
	return &pb.Vessel{
		Id:        v.ID,
		Capacity:  v.Capacity,
		MaxWeight: v.MaxWeight,
		Name:      v.Name,
		Available: v.Available,
		OwnerId:   v.OwnerID,
	}
}
func MarshalSpecification(v *pb.Specification) *Specification {
	return &Specification{
		Capacity:  v.Capacity,
		MaxWeight: v.MaxWeight,
	}
}
func UnmarshalSpecification(v *Specification) *pb.Specification {
	return &pb.Specification{
		Capacity:  v.Capacity,
		MaxWeight: v.MaxWeight,
	}
}

//Repository - interface of what our datastore should be
type repository interface {
	FindAvailable(ctx context.Context, spec *Specification) (*Vessel, error)
	Create(ctx context.Context, vessel *Vessel) error
}

//VesselRepository - the actual repository implementig our interface methods
type VesselRepository struct {
	collection *mongo.Collection
}

//NewVesselRepository - creates an empty repository
func NewVesselRepository(collection *mongo.Collection) *VesselRepository {
	return &VesselRepository{collection}
}

//FindAvailable - checks a specification against a slice of vessels
//if the capacity and max_weight are below the vessel's then we return that vessel
func (repo *VesselRepository) FindAvailable(ctx context.Context, spec *Specification) (*Vessel, error) {
	vessel := &Vessel{}
	//filter := bson.M{"capacity": bson.M{"$lte": spec.Capacity}, "max_weight": bson.M{"$lte": spec.MaxWeight}}

	if err := repo.collection.FindOne(ctx, bson.M{"name": "Boaty McBoatface"}).Decode(vessel); err != nil {
		return nil, err
	}
	return vessel, nil

}

func (repo *VesselRepository) Create(ctx context.Context, vessel *Vessel) error {
	log.Println("creating a vessel")
	_, err := repo.collection.InsertOne(ctx, vessel)
	return err
}
