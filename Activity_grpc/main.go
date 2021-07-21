package main

import (
	"log"
	"net"
	"fmt"

	"google.golang.org/grpc"
	"GitHub/Activity_grpc/protoservices"
	"GitHub/Activity_grpc/services"
)

func main() {
	fmt.Println("hello")
	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	UserServices := services.UserServices{}
	protoservices.RegisterUserServiceServer(grpcServer, &registerServices)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}