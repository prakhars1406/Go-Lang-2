package services

import (
	"context"
	"fmt"
	"log"

	"GitHub/Grpc_Users_Poc/client/protoservices"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
)

func GetUser(c protoservices.GetUserServiceClient) {
	fmt.Println()
	md := metadata.Pairs("token", "valid-tokennnn")
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	response, err := c.GetUser(ctx, &protoservices.GetUserRequest{UserId: "12345"})
	if err != nil {
		log.Fatalf("Error when calling LoginUser: %s", err)
	}
	b, err := protojson.Marshal(response)
	if err != nil {
		log.Println("Error in marshalling data ", err)
	}
	log.Printf("Response from server: %s", string(b))
}
