/*
This file is part of Domain package.
This package deals with giving user accounts.
This package also provides features such as get all users accounts,get users account by id etc.
*/
package Domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	utils "github.com/Channels/Utils"
)

var (
	AccountsDetailDao accountsDetailDaoInterface
)

func init() {
	AccountsDetailDao = &accountsDetailDao{}
}

type accountsDetailDaoInterface interface {
	GetAccountsDetail(string) (*AccountsSubsciptionDetail, *utils.ApplicationError)
	getAllAccount() (*Accounts, error)
	getAccountsIdDetails(string) (*AccountsDetail, error)
}

type accountsDetailDao struct{}

/*
Returns users balance and subscriprtion details
Parameters Required
	- accountID : Account Id for user whose details need to be returned

Return Parameters
	- AccountsSubsciptionDetail : Pointer to AccountsSubsciptionDetail struct
	- utils.ApplicationError : This returns error details in case of failure

Steps for function Flow
	1. Get user account based on his Id
	2. If user found then proceed else return error
	3. Get the pack details based on pack Id
	4. If pack found then proceed else return error
	5. Return Pointer to AccountsSubsciptionDetail struct

Note:This function is exported function so can be used by other packages.
*/
func (u *accountsDetailDao) GetAccountsDetail(accountID string) (*AccountsSubsciptionDetail, *utils.ApplicationError) {
	if len(accountID) == 0 {
		return nil, &utils.ApplicationError{Message: "Account Id cannot be nil", StatusCode: http.StatusNotFound, Code: "not_found"}
	}
	accountsDetail, err := AccountsDetailDao.getAccountsIdDetails(accountID)
	if err != nil {
		if err.Error() == "not_found" {
			return nil, &utils.ApplicationError{Message: fmt.Sprintf("Account Id %v does not exists,please subscribe to use the service", accountID), StatusCode: http.StatusNotFound, Code: "accId_not_found"}
		} else {
			return nil, &utils.ApplicationError{Message: fmt.Sprintf("Error in accessing db"), StatusCode: http.StatusInternalServerError, Code: "db_error"}
		}
	}
	packs, err := PacksDetailDao.getPacksIdDetails(accountsDetail.PackId)
	if err != nil {
		if err.Error() == "not_found" {
			return nil, &utils.ApplicationError{Message: fmt.Sprintf("Pack Id %v does not exists", packs.ID), StatusCode: http.StatusNotFound, Code: "not_found"}
		} else {
			return nil, &utils.ApplicationError{Message: fmt.Sprintf("Error in accessing db"), StatusCode: http.StatusInternalServerError, Code: "db_error"}
		}
	}
	return &AccountsSubsciptionDetail{AccountBalance: accountsDetail.Balance, CurrentSubscription: packs.Name}, nil
}

/*
Returns users account details
Parameters Required
	- accountID : Account Id for user whose details need to be returned

Return Parameters
	- AccountsDetail : Pointer to AccountsDetail struct
	- error : This returns error message in case of failure

Steps for function Flow
	1. Get all accounts from accounts.json
	2. Search for account id in array returned in previous step
	3. If account id found return account details else return error

Note:This function will be only used by domain service thats why its not exported.
*/
func (u *accountsDetailDao) getAccountsIdDetails(accountId string) (*AccountsDetail, error) {
	file, err := ioutil.ReadFile("./Data/accounts.json")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	var accountsDetail Accounts
	json.Unmarshal(file, &accountsDetail)
	if len(accountsDetail.Accounts) == 0 {
		return nil, errors.New("not_found")
	}
	for _, account := range accountsDetail.Accounts {
		if account.ID == accountId {
			return &account, nil
		}
	}
	return nil, errors.New("not_found")
}

/*
Returns all users accounts
Parameters Required
	- No param

Return Parameters
	- Accounts : Pointer to Accounts struct
	- error : This returns error message in case of failure

Steps for function Flow
	1. Get all accounts from accounts.json
	2. If no error is fetching accounts.json return pointer to Accounts struct else return error

Note:This function will be only used by domain service thats why its not exported.
*/
func (u *accountsDetailDao) getAllAccount() (*Accounts, error) {
	file, err := ioutil.ReadFile("./Data/accounts.json")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	var accountsDetail Accounts
	json.Unmarshal(file, &accountsDetail)
	return &accountsDetail, nil
}
