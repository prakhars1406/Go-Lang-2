package services

import (
	"GitHub/Activity_grpc/protoservices"
	context "context"
	"google.golang.org/protobuf/encoding/protojson"
)

type UserServices struct {
}

func (s *UserServices) AddUser(ctx context.Context, in *protoservices.RegisterRequest) (*protoservices.RegisterResponse, error) {
	log.Printf("Welcome: %s", in.Name)
	b, err := protojson.Marshal(in)
	if err != nil {
		log.Println("Error in marshalling data ", err)
	}
	log.Printf("Data from user is %s", string(b))
	return &protoservices.NewUserResponse{UserId: "12345"}, nil
	return registerResponse, nil
}
