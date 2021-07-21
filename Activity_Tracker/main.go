package main

import (
	"log"
	"net"
	"fmt"

	"google.golang.org/grpc"
	"GitHub/Activity_Tracker/protoservices"
	"GitHub/Activity_Tracker/services"
)

/*
Parameters Required
	- None
Return Parameters
	- None
Steps for function Flow
	1. This function will create a listener and creates grpc server
	2. Then it will register activity service and user service to grpc server.
	3. At the end it will serve the grpc server.
*/
func main() {
	fmt.Println("hello")
	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	userService := services.UserService{}
	protoservices.RegisterUserServiceServer(grpcServer, &userService)

	activityService := services.ActivityService{}
	protoservices.RegisterActivityServiceServer(grpcServer, &activityService)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

}