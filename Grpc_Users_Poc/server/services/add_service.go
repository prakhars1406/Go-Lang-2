package services

import (
	"GitHub/Grpc_Users_Poc/server/dao"
	"GitHub/Grpc_Users_Poc/server/protoservices"
	"fmt"
)

type AddServices struct {
}

func (r *AddServices) AddService(stream protoservices.AddService_AddServiceServer) error {

	fmt.Printf("AddServices function was invoked with a streaming request\n")
	getAddServicerResponse, err := dao.AddServiceDao.AddServices(stream)
	if err != nil {
		return err
	}
	return stream.SendAndClose(&protoservices.AddServiceResponse{
		Message: getAddServicerResponse,
	})
}
