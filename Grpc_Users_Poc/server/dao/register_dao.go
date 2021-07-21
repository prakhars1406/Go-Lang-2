package dao

import (
	"GitHub/Grpc_Users_Poc/server/protoservices"
	"log"

	"google.golang.org/protobuf/encoding/protojson"
)

var (
	RegisterDao registerDaoInterface
)

func init() {
	RegisterDao = &registerDao{}
}

type registerDaoInterface interface {
	Register(*protoservices.RegisterRequest) (*protoservices.RegisterResponse, error)
}

type registerDao struct {
}

func (d *registerDao) Register(in *protoservices.RegisterRequest) (*protoservices.RegisterResponse, error) {
	log.Printf("Welcome: %s", in.Name)
	b, err := protojson.Marshal(in)
	if err != nil {
		log.Println("Error in marshalling data ", err)
	}
	log.Printf("Data from user is %s", string(b))
	return &protoservices.RegisterResponse{UserId: "12345"}, nil
}
