package main

import (
	"context"

	pb "github.com/zjjt/shippingGo/consignementService/proto/consignement"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//ce fichier gere les transaction avec la base de donnees choisi
//on defini ci dessous les modeles de donnees dans des structures qui les represente dans la base
//des fonctions Marshal et UnMarshal sont ecrites pour convertir des structures de donnees des
//protobuffs vers la bd choisie en l'occurence mongo db ici

//Consignement -model de document representant un consignement dans la bd
type Consignement struct {
	ID          string     `json:"id"`
	Weight      int32      `json:"weight"`
	Description string     `json:"description"`
	Containers  Containers `json:"containers"`
	VesselID    string     `json:"vessel_id"`
}

//Container -model de document representant un container dans la bd
type Container struct {
	ID         string `json:"id"`
	CustomerID string `json:"customer_id"`
	UserID     string `json:"user_id"`
}

//Containers - type alias representant un slice de pointeur Container
type Containers []*Container

//MarshalContainerCollection - converti un slice de protobuf *Container en Model *Container
func MarshalContainerCollection(containers []*pb.Container) []*Container {
	collection := make([]*Container, 0)
	for _, container := range containers {
		collection = append(collection, MarshalContainer(container))
	}
	return collection
}

//UnmarshalContainerCollection -converti un slice de  Model *Container en protobuf *Container
func UnmarshalContainerCollection(containers []*Container) []*pb.Container {
	collection := make([]*pb.Container, 0)
	for _, container := range containers {
		collection = append(collection, UnmarshalContainer(container))
	}
	return collection
}

//MarshalConsignementCollection -converti un slice de protobuf *Consignement en Model *Consignement
func MarshalConsignementCollection(consignements []*pb.Consignement) []*Consignement {
	collection := make([]*Consignement, 0)
	for _, consignement := range consignements {
		collection = append(collection, MarshalConsignement(consignement))
	}
	return collection
}

//UnmarshalConsignementCollection -converti un slice de Model *Consignement en protobuf *Consignement
func UnmarshalConsignementCollection(consignements []*Consignement) []*pb.Consignement {
	collection := make([]*pb.Consignement, 0)
	for _, consignement := range consignements {
		collection = append(collection, UnmarshalConsignement(consignement))
	}
	return collection
}

//MarshalContainer -converti un protobuf *Container en un Model *Container
func MarshalContainer(container *pb.Container) *Container {
	return &Container{
		ID:         container.Id,
		CustomerID: container.CustomerId,
		UserID:     container.UserId,
	}
}

//UnmarshalContainer -converti un Model *Container en un protobuf *Container
func UnmarshalContainer(container *Container) *pb.Container {
	return &pb.Container{
		Id:         container.ID,
		CustomerId: container.CustomerID,
		UserId:     container.UserID,
	}
}

//MarshalConsignement -converti un protobuf *Consignement en un Model *Consignement
func MarshalConsignement(consignement *pb.Consignement) *Consignement {
	return &Consignement{
		ID:          consignement.Id,
		Weight:      consignement.Weight,
		Description: consignement.Description,
		Containers:  MarshalContainerCollection(consignement.Containers),
		VesselID:    consignement.VesselId,
	}
}

//UnmarshalConsignement -converti un Model *Consignement en un protobuf *Consignement
func UnmarshalConsignement(consignement *Consignement) *pb.Consignement {
	return &pb.Consignement{
		Id:          consignement.ID,
		Weight:      consignement.Weight,
		Description: consignement.Description,
		Containers:  UnmarshalContainerCollection(consignement.Containers),
		VesselId:    consignement.VesselID,
	}
}

//Here we create our repository interface
type repository interface {
	Create(ctx context.Context, consignement *Consignement) error
	GetAll(ctx context.Context) ([]*Consignement, error)
}

//Here we instanciate our database handler which will implement the repository interface above
type MongoRepository struct {
	collection *mongo.Collection
}

//NewRepository creates a new mongo dbhandler and returns a pointer to it
func NewRepository(collection *mongo.Collection) *MongoRepository {
	return &MongoRepository{collection: collection}
}

//Create -insert one *Consignement model document into the db
func (repo MongoRepository) Create(ctx context.Context, consignement *Consignement) error {
	_, err := repo.collection.InsertOne(ctx, consignement)
	return err

}
func (repo *MongoRepository) GetAll(ctx context.Context) ([]*Consignement, error) {
	cur, err := repo.collection.Find(ctx, bson.M{})
	var consignements []*Consignement
	for cur.Next(ctx) {
		var consignement *Consignement
		if err := cur.Decode(&consignement); err != nil {
			return nil, err
		}
		consignements = append(consignements, consignement)
	}
	return consignements, err
}
