package main

import (
	"github.com/jinzhu/gorm"
	pb "github.com/zjjt/shippingGo/userService/proto/user"
)

type User struct {
	Id       string
	Name     string
	Company  string
	Email    string
	Password string
}

func MarshalUser(u *pb.User) *User {
	return &User{
		Id:       u.Id,
		Name:     u.Name,
		Company:  u.Company,
		Email:    u.Email,
		Password: u.Password,
	}
}
func UnmarshalUser(u *User) *pb.User {
	return &pb.User{
		Id:       u.Id,
		Name:     u.Name,
		Company:  u.Company,
		Email:    u.Email,
		Password: u.Password,
	}
}
func MarshalUserCollection(u []*pb.User) []*User {
	var collection []*User
	for _, v := range u {
		collection = append(collection, MarshalUser(v))
	}
	return collection
}
func UnmarshalUserCollection(u []*User) []*pb.User {
	var collection []*pb.User
	for _, v := range u {
		collection = append(collection, UnmarshalUser(v))
	}
	return collection
}

type repository interface {
	Create(user *pb.User) error
	Get(id string) (*pb.User, error)
	GetAll() ([]*pb.User, error)
	GetByEmail(email string) (*pb.User, error)
}
type UserRepository struct {
	db *gorm.DB
}

func newUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}
func (repo *UserRepository) Create(user *pb.User) error {
	if err := repo.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}
func (repo *UserRepository) Get(id string) (*pb.User, error) {
	var user *pb.User
	user.Id = id
	if err := repo.db.First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
func (repo *UserRepository) GetAll() ([]*pb.User, error) {
	var users []*pb.User
	if err := repo.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil

}
func (repo *UserRepository) GetByEmail(email string) (*pb.User, error) {
	var user *pb.User
	if err := repo.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
