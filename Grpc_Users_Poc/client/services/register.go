package services

import (
	"context"
	"fmt"
	"log"

	"GitHub/Grpc_Users_Poc/client/protoservices"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
)

func Register(c protoservices.RegisterServiceClient) {
	fmt.Println()
	md := metadata.Pairs("token", "valid-token")
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	response, err := c.Register(ctx, &protoservices.RegisterRequest{Name: "Prakhar", Email: "prakhars123", Phone: "1234567890",
		Address: "house no 111,park street,GkP"})
	if err != nil {
		log.Fatalf("Error when calling RegisterNewUser: %s", err)
	}
	b, err := protojson.Marshal(response)
	if err != nil {
		log.Println("Error in marshalling data ", err)
	}
	log.Printf("Response from server: %s", string(b))
}
