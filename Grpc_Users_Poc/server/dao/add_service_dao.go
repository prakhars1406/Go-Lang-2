package dao

import (
	"GitHub/Grpc_Users_Poc/server/protoservices"
	"io"
	"log"
)

var (
	AddServiceDao addServiceInterface
)

func init() {
	AddServiceDao = &addServiceDao{}
}

type addServiceInterface interface {
	AddServices(protoservices.AddService_AddServiceServer) (string, error)
}

type addServiceDao struct {
}

func (d *addServiceDao) AddServices(stream protoservices.AddService_AddServiceServer) (string, error) {
	log.Println("Receiving service name")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// we have finished reading the client stream
			return "Service added successfully", nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		serviceName := req.GetServiceName()
		log.Println(serviceName)
	}
}
