/*
This file is part of Domain package.
This package deals with recharging user account and updating users remaing balance.
*/
package Domain

import (
	"encoding/json"
	"net/http"
	"time"

	utils "github.com/Channels/Utils"
)

var (
	RechargeAmountDao rechargeAmountDaoInterface
)

func init() {
	RechargeAmountDao = &rechargeAmountDao{}
}

type rechargeAmountDaoInterface interface {
	RechargeAmount(int, string) (float32, *utils.ApplicationError)
	updateBalance(float32, string) (bool, error)
}

type rechargeAmountDao struct{}

/*
Recharge users account and add balance to it
Parameters Required
	- amount : user's email id
	- accountID : user's account id

Return Parameters
	- float32 : Update account balance.
	- utils.ApplicationError : This required error details in case of failure

Steps for function Flow
	1. Search account id present in accounts.json or not
	2. If not present return error else proceed
	3. Update user account balance
	4. If success return new account balance else return error

Note:This function is exported function so can be used by other packages.
*/
func (u *rechargeAmountDao) RechargeAmount(amount int, accountId string) (float32, *utils.ApplicationError) {
	accountsDetail, err := AccountsDetailDao.getAccountsIdDetails(accountId)
	if err != nil {
		if err.Error() == "not_found" {
			return 100, &utils.ApplicationError{Message: "Account Id not subscribed", StatusCode: http.StatusNotFound, Code: "accId_not_found"}
		} else {
			return 0, &utils.ApplicationError{Message: "Error in accessing db", StatusCode: http.StatusInternalServerError, Code: "db_error"}
		}
	}
	success, err := RechargeAmountDao.updateBalance(accountsDetail.Balance+float32(amount), accountId)
	if err != nil {
		return 0, &utils.ApplicationError{Message: "Error in updating db", StatusCode: http.StatusInternalServerError, Code: "db_error"}
	}
	if success == true {
		return accountsDetail.Balance + float32(amount), nil
	}
	return 0, &utils.ApplicationError{Message: "Unknown error occured", StatusCode: http.StatusInternalServerError, Code: "unknown_error"}
}

/*
Updating user account balance
Parameters Required
	- newBalance : New balance which need to be udpated in accounts.json
	- accountID : user's account id

Return Parameters
	- bool : This tells function was success or not.
	- error : This returns error message in case of failure

Steps for function Flow
	1. Get all the records present in accounts.json
	2. Search for user's account id
	3. If account id present then update the balance.
	4. Return the success status as true or false.If error then return eroor.

Note:This function will be only used by domain service thats why its not exported.
*/
func (u *rechargeAmountDao) updateBalance(newBalance float32, accountID string) (bool, error) {
	existingAcc, err := AccountsDetailDao.getAllAccount()
	if err != nil {
		return false, err
	}
	success := true
	for i, _ := range existingAcc.Accounts {
		if existingAcc.Accounts[i].ID == accountID {
			existingAcc.Accounts[i].Balance = newBalance
			currentTime := time.Now()
			existingAcc.Accounts[i].RechargeDate = currentTime.Format("01-02-2006")
		}
	}
	byteArray, err := json.Marshal(existingAcc)
	if err != nil {
		return false, err
	}
	err = SubscribeDao.beginningUpdate(byteArray)
	if err != nil {
		return false, err
	}
	return success, nil
}
