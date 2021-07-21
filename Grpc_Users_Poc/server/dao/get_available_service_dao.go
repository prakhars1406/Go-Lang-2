package dao

import (
	"GitHub/Grpc_Users_Poc/server/protoservices"
	"strconv"
)

var (
	GetAavailableServiceDao getAavailableServiceInterface
)

func init() {
	GetAavailableServiceDao = &getAavailableServiceDao{}
}

type getAavailableServiceInterface interface {
	GetService(*protoservices.GetAavailableServiceRequest) ([]string, error)
}

type getAavailableServiceDao struct {
}

func (d *getAavailableServiceDao) GetService(in *protoservices.GetAavailableServiceRequest) ([]string, error) {
	result := make([]string, 5)
	for i := 0; i < 5; i++ {
		result[i] = "Service_" + strconv.Itoa(i)
	}
	return result, nil
}
