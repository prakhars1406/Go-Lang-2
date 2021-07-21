package dao

import (
	"GitHub/Grpc_Users_Poc/server/protoservices"
	"log"
)

var (
	GetuserDao getUserDaoInterface
)

func init() {
	GetuserDao = &getUserDao{}
}

type getUserDaoInterface interface {
	GetUser(in *protoservices.GetUserRequest) (*protoservices.GetUserResponse, error)
}

type getUserDao struct {
}

func (d *getUserDao) GetUser(in *protoservices.GetUserRequest) (*protoservices.GetUserResponse, error) {
	log.Printf("Welcome: user with Id:%s", in.UserId)
	return &protoservices.GetUserResponse{UserId: "1234", Name: "Prakhar", Email: "prakhars123", Phone: "1234567890",
		Address: "house no 111,park street,GkP"}, nil
}
