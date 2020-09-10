/*
This file is part of Domain package.
This package deals with user subsciption feature
This package can also provides feaures such as subscribe new users to base pack,subcribe existing user etc
*/
package Domain

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	utils "github.com/Channels/Utils"
)

var (
	SubscribeDao subscribeDaoInterface
)

func init() {
	SubscribeDao = &subscribeDao{}
}

type subscribeDaoInterface interface {
	SubscribeToBase(string, string, int) (bool, bool, *utils.ApplicationError)
	beginningUpdate([]byte) error
	subcribeNewUser(float32, string, string, int) error
	subscribeExistingUserData(float32, string, string, string, []int, string, string) (bool, error)
}

type subscribeDao struct{}

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
	1. Search account id present in accounts.json or not
	2. If not present then subscibe as new user,else if user present subscribe as existing user
	3. Return success and notify status if subscription is successful else return error details

Note:This function is exported function so can be used by other packages.
*/
func (u *subscribeDao) SubscribeToBase(accountID string, pack string, months int) (bool, bool, *utils.ApplicationError) {
	if len(accountID) == 0 {
		return false, false, &utils.ApplicationError{Message: "Account Id cannot be nil", StatusCode: http.StatusNotFound, Code: "not_found"}
	}
	accountsDetail, err := AccountsDetailDao.getAccountsIdDetails(accountID)
	if err != nil {
		if err.Error() == "not_found" {
			if strings.ToLower(pack) == "s" {
				if months > 2 {
					return false, false, &utils.ApplicationError{Message: "Not sufficient balance", StatusCode: http.StatusBadRequest, Code: "low_balance"}
				} else if months <= 2 {
					err := SubscribeDao.subcribeNewUser(float32(100-months*50), accountID, pack, months)
					if err != nil {
						return false, false, &utils.ApplicationError{Message: "Error in accessing db", StatusCode: http.StatusInternalServerError, Code: "db_error"}
					} else {
						return true, false, nil
					}
				}
			} else if strings.ToLower(pack) == "g" {
				if months > 1 {
					return false, false, &utils.ApplicationError{Message: "Not sufficient balance", StatusCode: http.StatusBadRequest, Code: "low_balance"}
				} else if months == 1 {
					err = SubscribeDao.subcribeNewUser(float32(100-months*100), accountID, pack, months)
					if err != nil {
						return false, false, &utils.ApplicationError{Message: "Error in accessing db", StatusCode: http.StatusInternalServerError, Code: "db_error"}
					} else {
						return true, false, nil
					}
				}
			}
		} else {
			return false, false, &utils.ApplicationError{Message: "Error in accessing db", StatusCode: http.StatusInternalServerError, Code: "db_error"}
		}
	} else {
		if strings.ToLower(pack) == "s" {
			discount := 0
			if months >= 3 {
				discount = ((months * 50) / 100) * 10
			}
			if float32(months*50-discount) <= accountsDetail.Balance {
				notify, err := SubscribeDao.subscribeExistingUserData(accountsDetail.Balance-float32(months*50-discount), accountID, pack, accountsDetail.RechargeDate, accountsDetail.ExtraChannels, accountsDetail.Email, accountsDetail.Phone)
				if err != nil {
					return false, false, &utils.ApplicationError{Message: "Error in accessing db", StatusCode: http.StatusInternalServerError, Code: "db_error"}
				} else {
					if notify == false {
						return true, false, nil
					}
				}
			} else {
				return false, false, &utils.ApplicationError{Message: "Not sufficient balance", StatusCode: http.StatusBadRequest, Code: "low_balance"}
			}
		} else if strings.ToLower(pack) == "g" {
			discount := 0
			if months >= 3 {
				discount = ((months * 100) / 100) * 10
			}
			if float32(months*100-discount) <= accountsDetail.Balance {
				notify, err := SubscribeDao.subscribeExistingUserData(accountsDetail.Balance-float32(months*100-discount), accountID, pack, accountsDetail.RechargeDate, accountsDetail.ExtraChannels, accountsDetail.Email, accountsDetail.Phone)
				if err != nil {
					return false, false, &utils.ApplicationError{Message: "Error in accessing db", StatusCode: http.StatusInternalServerError, Code: "db_error"}
				} else {
					if notify == false {
						return true, false, nil
					}
				}
			} else {
				return false, false, &utils.ApplicationError{Message: "Not sufficient balance", StatusCode: http.StatusBadRequest, Code: "low_balance"}
			}
		}
	}
	return true, true, nil
}

/*
Subscribes existing user to base packs
Parameters Required
	- remainingBalance : Balance remaining in users accoutn after subsription.
	- accountID : user's account id.
	- pack : Pack Id for which user is subscibing.
	- rechargeDate : Last recharge date for user
	- existingChannels : Array of channel id already present in user account
	- email : Email of user account
	- phone : Phone number of user account

Return Parameters
	- bool : Success Status.
	- error : Return error message in case of error

Steps for function Flow
	1. Get all accounts present in accounts.json
	2. Search for account id in data returned in previous step
	3. Update user details and balance
	4. If evert thing is success return success status else return failure message

Note:This function will be only used by domain service thats why its not exported.
*/
func (u *subscribeDao) subscribeExistingUserData(remainingBalance float32, accountID string, pack string, rechargeDate string, existingChannels []int, email string, phone string) (bool, error) {
	PackId := 0
	if strings.ToLower(pack) == "s" {
		PackId = 1
	} else if strings.ToLower(pack) == "g" {
		PackId = 2
	}
	existingAcc, err := AccountsDetailDao.getAllAccount()
	if err != nil {
		return false, err
	}
	notify := false
	for i, _ := range existingAcc.Accounts {
		if existingAcc.Accounts[i].ID == accountID {
			existingAcc.Accounts[i].PackId = PackId
			existingAcc.Accounts[i].Balance = remainingBalance
			currentTime := time.Now()
			existingAcc.Accounts[i].SubscriptionDate = currentTime.Format("01-02-2006")
			if len(existingAcc.Accounts[i].Email) != 0 {
				notify = true
			}
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
	return notify, nil
}

/*
Subscribes new user to base packs
Parameters Required
	- Balance : Balance remaining in users accoutn after subsription.
	- accountID : user's account id.
	- pack : Pack Id for which user is subscibing.
	- months : Number of months for which user want to subscribe

Return Parameters
	- error : Return error message in case of error

Steps for function Flow
	1. Create new users data.
	2. Get all accounts present in accounts.json
	3. Append new users data to data returned in step 2.
	3. Update accounts.json file
	4. If evert thing is success return error as nil,else return error message.

Note:This function will be only used by domain service thats why its not exported.
*/
func (u *subscribeDao) subcribeNewUser(balance float32, accountID string, pack string, months int) error {
	var newaccount AccountsDetail
	PackId := 0
	if strings.ToLower(pack) == "s" {
		PackId = 1
	} else if strings.ToLower(pack) == "g" {
		PackId = 2
	}
	newaccount.ID = accountID
	newaccount.PackId = PackId
	newaccount.Balance = balance
	currentTime := time.Now()
	newaccount.SubscriptionDate = currentTime.Format("01-02-2006")
	newaccount.RechargeDate = currentTime.Format("01-02-2006")
	newaccount.ExtraChannels = []int{}
	newaccount.Email = ""
	newaccount.Phone = ""
	existingAcc, err := AccountsDetailDao.getAllAccount()
	if err != nil {
		return err
	}
	existingAcc.Accounts = append(existingAcc.Accounts, newaccount)
	byteArray, err := json.Marshal(existingAcc)
	if err != nil {
		return err
	}
	err = SubscribeDao.beginningUpdate(byteArray)
	if err != nil {
		return err
	}
	return nil
}

/*
Updae json file
Parameters Required
	- data : Data which needs to be updated in json file.

Return Parameters
	- error : Return error message in case of error

Steps for function Flow
	1. Open accounts.json file
	2. Update data in accounts.json file
	3. If evert thing is success return error as nil,else return error message.

Note:This function will be only used by domain service thats why its not exported.
*/
func (u *subscribeDao) beginningUpdate(data []byte) error {
	// Read Write Mode
	file, err := os.OpenFile("./Data/accounts.json", os.O_RDWR, 0644)

	if err != nil {
		return err
	}
	defer file.Close()
	s := ""
	_, err = file.WriteAt([]byte(s), 0) // Write at 0 beginning
	_, err = file.WriteAt(data, 0)      // Write at 0 beginning
	if err != nil {
		return err
	}
	return nil
}
