/*
This file is part of Services package.
This package deals with giving user info about packs and channels.
This package also provides feature to user to add new channels.
*/
package Services

import (
	domain "github.com/Channels/Domain"
	utils "github.com/Channels/Utils"
)

type packsService struct{}

var (
	PacksService packsServiceInterface
)

func init() {
	PacksService = &packsService{}
}

type packsServiceInterface interface {
	GetPacksAndChannels() (*domain.PacksAndChannels, *utils.ApplicationError)
	AddChannel(string, string) (float32, *utils.ApplicationError)
}

/*
Return the available packs and channels
Parameters Required
	-No Param

Return Parameters
	- PacksAndChannels : Pointer to PacksAndChannels struct
	- utils.ApplicationError : This returns error details in case of failure

Steps for function Flow
	1. This function will internall call domain function.
	2. Returns the value returned by domanin function

Note:This function is exported function so can be used by other packages.
*/
func (pack *packsService) GetPacksAndChannels() (*domain.PacksAndChannels, *utils.ApplicationError) {
	accountSubs, err := domain.PacksDetailDao.GetPacksAndChannels()
	if err != nil {
		return nil, err
	}
	return accountSubs, nil
}

/*
Add new channels to user account
Parameters Required
	- channel : new channel name to be added
	- accountId : account Id to which channel should be added

Return Parameters
	- int : balance remaining after adding the channel
	- utils.ApplicationError : This required error details in case of failure

Steps for function Flow
	1. This function will internall call domain function.
	2. Returns the value returned by domanin function

Note:This function is exported function so can be used by other packages.
*/
func (pack *packsService) AddChannel(channel string, accountId string) (float32, *utils.ApplicationError) {
	balance, err := domain.PacksDetailDao.AddChannel(channel, accountId)
	if err != nil {
		return 0, err
	}
	return balance, nil
}
