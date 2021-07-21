package services

import (
	"GitHub/Grpc_Users_Poc/server/dao"
	"GitHub/Grpc_Users_Poc/server/protoservices"
	"log"
	"time"
)

type GetAavailableServices struct {
}

func (r *GetAavailableServices) GetService(in *protoservices.GetAavailableServiceRequest, stream protoservices.GetAavailableService_GetServiceServer) error {
	log.Printf("Request from user:%s", in.Name)
	getServicerResponse, err := dao.GetAavailableServiceDao.GetService(in)
	if err != nil {
		return err
	}
	for _, v := range getServicerResponse {
		res := &protoservices.GetAavailableServiceResponse{
			ServiceName: v,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}
