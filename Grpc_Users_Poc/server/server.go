package main

import (
	"GitHub/Grpc_Users_Poc/server/protoservices"
	"GitHub/Grpc_Users_Poc/server/services"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

func main() {
	fmt.Println("Go grpc server!")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	registerServices := services.RegisterServices{}
	protoservices.RegisterRegisterServiceServer(grpcServer, &registerServices)

	getUserServices := services.GetUserServices{}
	protoservices.RegisterGetUserServiceServer(grpcServer, &getUserServices)

	getAvailableServices := services.GetAavailableServices{}
	protoservices.RegisterGetAavailableServiceServer(grpcServer, &getAvailableServices)

	addServices := services.AddServices{}
	protoservices.RegisterAddServiceServer(grpcServer, &addServices)

	checkUserServices := services.CheckUserServices{}
	protoservices.RegisterCheckUserServiceServer(grpcServer, &checkUserServices)
	go func() {
		fmt.Println("Starting Server...")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %s", err)
		}
	}()
	// Wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	// Block until a signal is received
	<-ch
	if err := lis.Close(); err != nil {
		log.Fatalf("Error on closing the listener : %v", err)
	}
	// Finally, we stop the server
	fmt.Println("Stopping the server")
	grpcServer.Stop()
	fmt.Println("End of Program")
}
func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, grpc.Errorf(codes.Unauthenticated, "missing context metadata")
	}
	if len(meta["token"]) != 1 {
		return nil, grpc.Errorf(codes.Unauthenticated, "invalid token")
	}
	if meta["token"][0] != "valid-token" {
		return nil, grpc.Errorf(codes.Unauthenticated, "invalid token")
	}

	return handler(ctx, req)
}
