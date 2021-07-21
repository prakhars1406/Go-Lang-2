package services

import (
	"log"
	"google.golang.org/protobuf/encoding/protojson"
	"GitHub/Grpc_Users_Poc/server/protoservices"
	context "context"
)

type RegisterServices struct {
}

func (s *RegisterServices) Register(ctx context.Context, in *protoservices.RegisterRequest) (*protoservices.RegisterResponse, error) {

	log.Printf("Welcome: %s", in.Name)
	b, err := protojson.Marshal(in)
	if err != nil {
		log.Println("Error in marshalling data ", err)
	}
	log.Printf("Data from user is %s", string(b))
	return &protoservices.RegisterResponse{UserId: "12345"}, nil
}
