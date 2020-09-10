/*
This file is part of Domain package.
This package deals with updating user's email and phone number.
This package can also be extended to send notification to user.Currently this feature is not supported
*/
package Domain

import (
	"encoding/json"
	"errors"
	"net/http"

	utils "github.com/Channels/Utils"
)

var (
	NotificationDao notificationDaoInterface
)

func init() {
	NotificationDao = &notficationDao{}
}

type notificationDaoInterface interface {
	UpdateEmailAndPhone(string, string, string) (bool, *utils.ApplicationError)
	updateEmailAndPhone(string, string, string) (bool, error)
}

type notficationDao struct{}

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
	1. Search account id present in accounts.json or not
	2. If not present return error else proceed
	3. Call updateEmailAndPhone to update the email and phone number to accounts.json file
	4. Return the success status as true or false.If error then return error.

Note:This function is exported function so can be used by other packages.
*/
func (notify *notficationDao) UpdateEmailAndPhone(email string, phone string, accountId string) (bool, *utils.ApplicationError) {
	_, err := AccountsDetailDao.getAccountsIdDetails(accountId)
	if err != nil {
		if err.Error() == "not_found" {
			return false, &utils.ApplicationError{
				Message:    "Account Id not subscribed",
				StatusCode: http.StatusNotFound,
				Code:       "accId_not_found",
			}
		} else {
			return false, &utils.ApplicationError{
				Message:    "Error in accessing db",
				StatusCode: http.StatusInternalServerError,
				Code:       "db_error",
			}
		}
	}
	success, err := NotificationDao.updateEmailAndPhone(email, phone, accountId)
	if err != nil {
		if err.Error() == "already_present" {
			return false, &utils.ApplicationError{
				Message:    "Email and phone already updated",
				StatusCode: http.StatusMethodNotAllowed,
				Code:       "already_present",
			}
		}
		return false, &utils.ApplicationError{
			Message:    "Error in updating db",
			StatusCode: http.StatusInternalServerError,
			Code:       "db_error",
		}
	}
	if success == true {
		return true, nil
	}
	return false, &utils.ApplicationError{Message: "Unknown error occured", StatusCode: http.StatusInternalServerError, Code: "unknown_error"}
}

/*
Updating user email and phone number in accounts.json file
Parameters Required
	- email : user's email id
	- phone : user's phone number
	- accountID : user's account id

Return Parameters
	- bool : This tells function was success or not.
	- error : This returns error message in case of failure

Steps for function Flow
	1. Get all the records present in accounts.json
	2. Search for user's account id
	3. If account id present then update the email and phone number.
	4. Update the json file with newly data updated data.
	5. Return the success status as true or false.If error then return error.

Note:This function will be only used by domain service thats why its not exported.
*/
func (notify *notficationDao) updateEmailAndPhone(email string, phone string, accountID string) (bool, error) {
	existingAcc, err := AccountsDetailDao.getAllAccount()
	if err != nil {
		return false, err
	}
	success := true
	for i, _ := range existingAcc.Accounts {
		if existingAcc.Accounts[i].ID == accountID {
			if len(existingAcc.Accounts[i].Email) != 0 {
				return false, errors.New("already_present")
			}
			existingAcc.Accounts[i].Email = email
			existingAcc.Accounts[i].Phone = phone
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
