/*
This file is part of Services package.
This package deals with giving user accounts.
This package also provides features such as get all users accounts,get users account by id etc.
*/
package Services

import (
	domain "github.com/Channels/Domain"
	utils "github.com/Channels/Utils"
)

type notificationService struct{}

var (
	NotificationService notificationServiceInterface
)

func init() {
	NotificationService = &notificationService{}
}

type notificationServiceInterface interface {
	UpdateEmailAndPhone(string, string, string) (bool, *utils.ApplicationError)
}

/*
Updating user email and phone number in accounts.json file
Parameters Required
	- email : user's email id
	- phone : user's phone number
	- accountID : user's account id

Return Parameters
	- bool : This tells function was success or not.
	- utils.ApplicationError : This required error details in case of failure

Steps for function Flow
	1. This function will internall call domain function.
	2. Returns the value returned by domanin function

Note:This function is exported function so can be used by other packages.
*/
func (notify *notificationService) UpdateEmailAndPhone(email string, phone string, accountId string) (bool, *utils.ApplicationError) {
	success, err := domain.NotificationDao.UpdateEmailAndPhone(email, phone, accountId)
	if err != nil {
		return success, err
	}
	return success, nil
}
