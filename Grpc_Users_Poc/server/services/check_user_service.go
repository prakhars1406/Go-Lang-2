package services

import (
	"GitHub/Grpc_Users_Poc/server/dao"
	"GitHub/Grpc_Users_Poc/server/protoservices"
	"fmt"
)

type CheckUserServices struct {
}

func (r *CheckUserServices) CheckUserService(stream protoservices.CheckUserService_CheckUserServiceServer) error {

	fmt.Println("Starting to do a CheckServices BiDi Streaming RPC...")
	err := dao.CheckUserServiceDao.CheckUserService(stream)
	if err != nil {
		return err
	}
	return nil
}
