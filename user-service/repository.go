package main

import (
	pb "../user-service/proto/user"
	// "github.com/jinzhu/gorm"
)

type Repository interface {
	GetAll() ([]*pb.User, error)
	Get(id string) (*pb.User, error)
	Create(user *pb.User) error
	GetByEmail(email string) (*pb.User, error)
}

type UserRepository struct {
	userList map[string]*pb.User
}

func (repo *UserRepository) New() error {
	repo.userList = make(map[string]*pb.User)
	return nil
}

func (repo *UserRepository) GetAll() ([]*pb.User, error) {
	var users []*pb.User

	for _, v := range repo.userList {
		users = append(users, v)
	}
	return users, nil
}

func (repo *UserRepository) Get(id string) (*pb.User, error) {
	var user *pb.User
	user.Id = id
	for _, v := range repo.userList {
		if v.Id == id {
			return user, nil
		}
	}
	return nil, nil
}

func (repo *UserRepository) GetByEmail(email string) (*pb.User, error) {
	return repo.userList[email], nil
}

func (repo *UserRepository) Create(user *pb.User) error {
	repo.userList[user.Email] = user
	return nil
}
