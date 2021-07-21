package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"GitHub/Grpc_Users_Poc/graphQl/graph/generated"
	"GitHub/Grpc_Users_Poc/graphQl/graph/model"
	"GitHub/Grpc_Users_Poc/graphQl/protoservices"
	"context"
	"errors"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
)

var conn *grpc.ClientConn

func init() {
	con, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	conn = con
}
func (r *mutationResolver) Register(ctx context.Context, input model.RegisterRequest) (*model.RegisterResponse, error) {
	register := model.RegisterResponse{}
	registerServiceClient := protoservices.NewRegisterServiceClient(conn)
	md := metadata.Pairs("token", "valid-token")
	ctx = metadata.NewOutgoingContext(ctx, md)
	response, err := registerServiceClient.Register(ctx, &protoservices.RegisterRequest{Name: "Prakhar", Email: "prakhars123", Phone: "1234567890",
		Address: "house no 111,park street,GkP"})
	if err != nil {
		log.Fatalf("Error when calling RegisterNewUser: %s", err)
		return &register, err
	}
	b, err := protojson.Marshal(response)
	if err != nil {
		log.Println("Error in marshalling data ", err)
		return &register, err
	}
	log.Printf("Response from server: %s", string(b))

	getUserServiceClient := protoservices.NewGetUserServiceClient(conn)
	response2, err := getUserServiceClient.GetUser(ctx, &protoservices.GetUserRequest{UserId: response.GetUserId()})
	if err != nil {
		log.Fatalf("Error when calling LoginUser: %s", err)
		return &register, err
	}
	b, err = protojson.Marshal(response2)
	if err != nil {
		log.Println("Error in marshalling data ", err)
		return &register, err
	}
	log.Printf("Response from server: %s", string(b))
	register.UserID = response2.GetUserId()
	register.Name = response2.GetName()
	register.Email = response2.GetEmail()
	register.Phone = response2.GetPhone()
	register.Address = response2.GetAddress()
	fmt.Println(register)
	return &register, nil
}

func (r *queryResolver) GetUser(ctx context.Context, id string) (*model.RegisterResponse, error) {
	register := model.RegisterResponse{}
	if id == "1234" {
		md := metadata.Pairs("token", "valid-token")
		ctx = metadata.NewOutgoingContext(ctx, md)
		getUserServiceClient := protoservices.NewGetUserServiceClient(conn)
		response2, err := getUserServiceClient.GetUser(ctx, &protoservices.GetUserRequest{UserId: id})
		if err != nil {
			log.Fatalf("Error when calling LoginUser: %s", err)
			return &register, err
		}
		b, err := protojson.Marshal(response2)
		if err != nil {
			log.Println("Error in marshalling data ", err)
			return &register, err
		}
		log.Printf("Response from server: %s", string(b))
		register.UserID = response2.GetUserId()
		register.Name = response2.GetName()
		register.Email = response2.GetEmail()
		register.Phone = response2.GetPhone()
		register.Address = response2.GetAddress()
		fmt.Println(register)
		return &register, nil
	}
	return &register, errors.New("id not matching")
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
