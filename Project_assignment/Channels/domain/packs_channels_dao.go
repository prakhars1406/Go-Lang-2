/*
This file is part of Domain package.
This package deals with giving user info about packs and channels.
This package also provides feature to user to add new channels.
*/
package Domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	utils "github.com/Channels/Utils"
)

var (
	PacksDetailDao packsDetailDaoInterface
)

func init() {
	PacksDetailDao = &packsDetailDao{}
}

type packsDetailDaoInterface interface {
	GetPacksAndChannels() (*PacksAndChannels, *utils.ApplicationError)
	AddChannel(string, string) (float32, *utils.ApplicationError)
	getPacksIdDetails(int) (*PacksDetails, error)
	addchannel(int, int, string) (bool, error)
	getPacks() (*Packs, error)
	getChannels() (*Channels, error)
}

type packsDetailDao struct{}

/*
Return the available packs and channels
Parameters Required
	-No Param

Return Parameters
	- PacksAndChannels : Pointer to PacksAndChannels struct
	- utils.ApplicationError : This returns error details in case of failure

Steps for function Flow
	1. Get all the packs present in packs.json
	2. Get all the channels present in channels.json
	3. Return the PacksAndChannels struct pointer.

Note:This function is exported function so can be used by other packages.
*/
func (pack *packsDetailDao) GetPacksAndChannels() (*PacksAndChannels, *utils.ApplicationError) {
	packs, err := PacksDetailDao.getPacks()
	if err != nil {
		return nil, &utils.ApplicationError{Message: fmt.Sprintf("Unable to get packs from DB"), StatusCode: http.StatusInternalServerError, Code: "db_error"}
	}
	channels, err := PacksDetailDao.getChannels()
	if err != nil {
		return nil, &utils.ApplicationError{Message: fmt.Sprintf("Unable to get channles from DB"), StatusCode: http.StatusInternalServerError, Code: "db_error"}
	}
	return &PacksAndChannels{Packs: packs.Packs, Channels: channels.Channels}, nil
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
	1. Get all the channels present in channels.json
	2. Find channel is present or not,if not present return error
	3. If channel is present fetch channel id and channel price
	4. Get account details of user.
	5. Check channel already present or not.
	6. If channel present return error.
	7. If channel not present check for user balance.
	8. If user has sufficient balance then add channel to users account.
	9. If channel is added successfully then return balance else return error.

Note:This function is exported function so can be used by other packages.
*/
func (pack *packsDetailDao) AddChannel(channel string, accountId string) (float32, *utils.ApplicationError) {
	channels, err := PacksDetailDao.getChannels()
	if err != nil {
		return 0, &utils.ApplicationError{Message: fmt.Sprintf("Unable to get channles from DB"), StatusCode: http.StatusInternalServerError, Code: "db_error"}
	}
	channelPrice := 0
	channelId := 0
	for _, channelValue := range channels.Channels {
		if strings.ToLower(channelValue.Name) == channel {
			channelPrice = channelValue.Price
			channelId = channelValue.ID
		}
	}
	if channelPrice == 0 {
		return 0, &utils.ApplicationError{Message: fmt.Sprintf("Invalid channel selected"), StatusCode: http.StatusBadRequest, Code: "channel_not_found"}
	}
	accountsDetail, err := AccountsDetailDao.getAccountsIdDetails(accountId)
	if err != nil {
		if err.Error() == "not_found" {
			return 0, &utils.ApplicationError{Message: "Account Id not found", StatusCode: http.StatusNotFound, Code: "not_found"}
		} else {
			return 0, &utils.ApplicationError{Message: "Error in accessing db", StatusCode: http.StatusInternalServerError, Code: "db_error"}
		}
	}
	if accountsDetail.PackId == 2 {
		return 0, &utils.ApplicationError{Message: "Channel already present", StatusCode: http.StatusMethodNotAllowed, Code: "already_present"}
	} else if accountsDetail.PackId == 1 {
		if channelId == 10 || channelId == 11 || channelId == 12 {
			return 0, &utils.ApplicationError{Message: "Channel already present", StatusCode: http.StatusMethodNotAllowed, Code: "already_present"}
		} else {
			_, present := utils.Seach(accountsDetail.ExtraChannels, channelId)
			if present == true {
				return 0, &utils.ApplicationError{Message: "Channel already present", StatusCode: http.StatusMethodNotAllowed, Code: "already_present"}
			}
		}
	}
	if accountsDetail.Balance >= float32(channelPrice) {
		success, err := PacksDetailDao.addchannel(channelId, channelPrice, accountId)
		if err != nil {
			return 0, &utils.ApplicationError{Message: "Error in accessing db", StatusCode: http.StatusInternalServerError, Code: "db_error"}
		} else if success == true {
			return accountsDetail.Balance - float32(channelPrice), nil
		}
	}
	return 0, &utils.ApplicationError{Message: "Unknown error", StatusCode: http.StatusInternalServerError, Code: "unknown_error"}
}

/*
This will add channel to users account
Parameters Required
	- channelId : Channel id to be added to account.
	- channlePrice : Price of channel to be added to account
	- accountID : User account Id to which channel needs to be added

Return Parameters
	- bool : This return true if channel added successfully,else return false
	- error : This returns error message in case of failure

Steps for function Flow
	1. Get all accounts present in accounts.json
	2. Search for account id of user
	3. If account id found add new channels and update new balance
	4. If channel added successfully return true else return false

Note:This function will be only used by domain service thats why its not exported.
*/
func (pack *packsDetailDao) addchannel(channelId int, channlePrice int, accountID string) (bool, error) {
	existingAcc, err := AccountsDetailDao.getAllAccount()
	if err != nil {
		return false, err
	}
	success := true
	for i, _ := range existingAcc.Accounts {
		if existingAcc.Accounts[i].ID == accountID {
			existingAcc.Accounts[i].Balance -= float32(channlePrice)
			existingAcc.Accounts[i].ExtraChannels = append(existingAcc.Accounts[i].ExtraChannels, channelId)
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

/*
Return the available packs
Parameters Required
	-No Param

Return Parameters
	- *Packs : pointer to pack struct.
	- error : This returns error message in case of failure

Steps for function Flow
	1. Read packs.json file
	2. Convert byte array to packs struct
	3. Return pointer to pack struct.

Note:This function will be only used by domain service thats why its not exported.
*/
func (pack *packsDetailDao) getPacks() (*Packs, error) {
	file, err := ioutil.ReadFile("./Data/Packs.json")
	if err != nil {
		return nil, err
	}
	var packs Packs
	json.Unmarshal(file, &packs)
	return &packs, nil
}

/*
Return the available channels
Parameters Required
	-No Param

Return Parameters
	- *Channels : pointer to Channels struct.
	- error : This returns error message in case of failure

Steps for function Flow
	1. Read channels.json file
	2. Convert byte array to channels struct
	3. Return pointer to channels struct.

Note:This function will be only used by domain service thats why its not exported.
*/
func (pack *packsDetailDao) getChannels() (*Channels, error) {
	file, err := ioutil.ReadFile("./Data/channels.json")
	if err != nil {
		return nil, err
	}
	var chanels Channels
	json.Unmarshal(file, &chanels)
	return &chanels, nil
}

/*
Return the packs for specific pack Id
Parameters Required
	- packId : pack Id for which details is required.

Return Parameters
	- *PacksDetails : pointer to PacksDetails struct.
	- error : This returns error message in case of failure

Steps for function Flow
	1. Read packs.json file
	2. Search for pack Id
	3. If pack Id found return pointer to PacksDetails struct.

Note:This function will be only used by domain service thats why its not exported..
*/
func (pack *packsDetailDao) getPacksIdDetails(packId int) (*PacksDetails, error) {
	file, err := ioutil.ReadFile("./Data/Packs.json")
	if err != nil {
		return nil, err
	}
	var packsDetails Packs
	json.Unmarshal(file, &packsDetails)
	for _, pack := range packsDetails.Packs {
		if pack.ID == packId {
			return &pack, nil
		}
	}
	return nil, errors.New("not_found")
}
