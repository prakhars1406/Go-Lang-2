package services

import (
	"context"
	"fmt"
	"log"

	"GitHub/Grpc_Users_Poc/client/protoservices"
)

func AddService(c protoservices.AddServiceClient) {
	fmt.Println()
	// md := metadata.Pairs("token", "cccvalid-tokennn")
	// ctx := metadata.NewOutgoingContext(context.Background(), md)
	stream, err := c.AddService(context.Background())
	if err != nil {
		log.Fatalf("Error while opening stream: %v", err)
	}

	services := []string{"Service_0", "Service_2", "Service_4"}

	for _, service := range services {
		fmt.Printf("Sending service: %s\n", service)
		stream.Send(&protoservices.AddServiceRequest{
			ServiceName: service,
		})
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response: %v", err)
	}

	fmt.Printf("The message is: %v\n", res.GetMessage())
}
