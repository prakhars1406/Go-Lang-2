package dao

import (
	"GitHub/Grpc_Users_Poc/server/protoservices"
	"fmt"
	"io"
	"log"
)

var (
	CheckUserServiceDao checkUserServiceInterface
)

func init() {
	CheckUserServiceDao = &checkUserServiceDao{}
}

type checkUserServiceInterface interface {
	CheckUserService(protoservices.CheckUserService_CheckUserServiceServer) error
}

type checkUserServiceDao struct {
}

func (d *checkUserServiceDao) CheckUserService(stream protoservices.CheckUserService_CheckUserServiceServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}
		serviceName := req.GetServiceName()
		fmt.Printf("Received a new serviceName of...: %s\n", serviceName)
		if serviceName == "Service_0" || serviceName == "Service_2" || serviceName == "Service_4" {
			sendErr := stream.Send(&protoservices.CheckUserServiceResponse{
				Message: serviceName + " already subcribed",
			})
			if sendErr != nil {
				log.Fatalf("Error while sending data to client: %v", sendErr)
				return sendErr
			}
		} else {
			sendErr := stream.Send(&protoservices.CheckUserServiceResponse{
				Message: serviceName + " not subcribed",
			})
			if sendErr != nil {
				log.Fatalf("Error while sending data to client: %v", sendErr)
				return sendErr
			}
		}
	}
}
