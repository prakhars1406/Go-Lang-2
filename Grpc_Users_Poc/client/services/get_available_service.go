package services

import (
	"context"
	"fmt"
	"io"
	"log"

	"GitHub/Grpc_Users_Poc/client/protoservices"

	"google.golang.org/grpc/metadata"
)

func GetService(c protoservices.GetAavailableServiceClient) {
	fmt.Println()
	md := metadata.Pairs("token", "valid-tokennn")
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	req := &protoservices.GetAavailableServiceRequest{
		Name: "Get all services",
	}

	resStream, err := c.GetService(ctx, req)
	if err != nil {
		log.Fatalf("error while calling GetServices RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("Response from GetServices: %v", msg.GetServiceName())
	}
}
