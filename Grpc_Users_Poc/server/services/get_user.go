package services

import (
	"GitHub/Grpc_Users_Poc/server/dao"
	"GitHub/Grpc_Users_Poc/server/protoservices"
	context "context"
)

type GetUserServices struct {
}

func (s *GetUserServices) GetUser(ctx context.Context, in *protoservices.GetUserRequest) (*protoservices.GetUserResponse, error) {
	registerResponse, err := dao.GetuserDao.GetUser(in)
	if err != nil {
		return nil, err
	}
	return registerResponse, nil
}
