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

type AccountsServ struct{}

var (
	AccountsService AccountsServiceInterface = AccountsServ{}
)

type AccountsServiceInterface interface {
	GetAccountsDetail(string) (*domain.AccountsSubsciptionDetail, *utils.ApplicationError)
	RechargeAmount(int, string) (float32, *utils.ApplicationError)
	SubscribeToBase(string, string, int) (bool, bool, *utils.ApplicationError)
}

/*
Returns users balance and subscriprtion details
Parameters Required
	- accountID : Account Id for user whose details need to be returned

Return Parameters
	- AccountsSubsciptionDetail : Pointer to AccountsSubsciptionDetail struct
	- utils.ApplicationError : This returns error details in case of failure

Steps for function Flow
	1. This function will internall call domain function.
	2. Returns the value returned by domanin function

Note:This function is exported function so can be used by other packages.
*/
func (account AccountsServ) GetAccountsDetail(accountID string) (*domain.AccountsSubsciptionDetail, *utils.ApplicationError) {
	accountSubs, err := domain.AccountsDetailDao.GetAccountsDetail(accountID)
	if err != nil {
		return nil, err
	}
	return accountSubs, nil
}

/*
Recharge users account and add balance to it
Parameters Required
	- amount : user's email id
	- accountID : user's account id

Return Parameters
	- float32 : Update account balance.
	- utils.ApplicationError : This required error details in case of failure

Steps for function Flow
	1. This function will internall call domain function.
	2. Returns the value returned by domanin function

Note:This function is exported function so can be used by other packages.
*/
func (account AccountsServ) RechargeAmount(amount int, accountId string) (float32, *utils.ApplicationError) {
	balance, err := domain.RechargeAmountDao.RechargeAmount(amount, accountId)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

/*
Subscribes user to base packs
Parameters Required
	- pack : Base pack in which user has to subscribe
	- accountID : user's account id
	- months : Number of months for which user has to subscribe

Return Parameters
	- bool : Success Status.
	- bool : Notify status
	- utils.ApplicationError : This required error details in case of failure

Steps for function Flow
	1. This function will internall call domain function.
	2. Returns the value returned by domanin function

Note:This function is exported function so can be used by other packages.
*/
func (account AccountsServ) SubscribeToBase(accountID string, pack string, months int) (bool, bool, *utils.ApplicationError) {
	success, notify, err := domain.SubscribeDao.SubscribeToBase(accountID, pack, months)
	if err != nil {
		return success, notify, err
	}
	return success, notify, nil
}
