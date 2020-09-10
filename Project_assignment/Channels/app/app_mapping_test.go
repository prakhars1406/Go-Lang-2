package app

import (
	"fmt"
	"net/http"
	"testing"

	domain "github.com/Channels/Domain"
	utils "github.com/Channels/Utils"
	"github.com/stretchr/testify/assert"
)

var (
	getAccountsDetail   func(accountId string) (*domain.AccountsSubsciptionDetail, *utils.ApplicationError)
	rechargeAmount      func(amount int, accountId string) (float32, *utils.ApplicationError)
	getPacksAndChannels func() (*domain.PacksAndChannels, *utils.ApplicationError)
	updateEmailAndPhone func(string, string, string) (bool, *utils.ApplicationError)
)

type AccountsService struct{}

func (s *AccountsService) GetAccountsDetail(accountId string) (*domain.AccountsSubsciptionDetail, *utils.ApplicationError) {
	return getAccountsDetail(accountId)
}
func (s *AccountsService) RechargeAmount(amount int, accountId string) (float32, *utils.ApplicationError) {
	return rechargeAmount(amount, accountId)
}
func (s *AccountsService) GetPacksAndChannels() (*domain.PacksAndChannels, *utils.ApplicationError) {
	return getPacksAndChannels()
}
func (s *AccountsService) UpdateEmailAndPhone(email string, phone string, accountId string) (bool, *utils.ApplicationError) {
	return updateEmailAndPhone(email, phone, accountId)
}

//Table driven test for GetAccountsDetails
var GetAccountsDetailsData = []struct {
	input          string
	accountDetails *domain.AccountsSubsciptionDetail
	err            *utils.ApplicationError
	output         int
}{
	{
		input:          "prakhars",
		accountDetails: &domain.AccountsSubsciptionDetail{AccountBalance: 100, CurrentSubscription: "Silver pack"},
		err:            nil,
		output:         100,
	},
	{
		input:          "prakhars",
		accountDetails: &domain.AccountsSubsciptionDetail{AccountBalance: 200, CurrentSubscription: "Gold pack"},
		err:            nil,
		output:         200,
	},
	{
		input:          "prakhars",
		accountDetails: &domain.AccountsSubsciptionDetail{AccountBalance: 350, CurrentSubscription: "Gold pack"},
		err:            nil,
		output:         350,
	},
}

func TestCallGetAccountsDetailsForExistingUser(t *testing.T) {
	accountsService := &AccountsService{}
	for _, v := range GetAccountsDetailsData {
		getAccountsDetail = func(accountId string) (*domain.AccountsSubsciptionDetail, *utils.ApplicationError) {
			return v.accountDetails, v.err
		}
		accountsDetails, _ := accountsService.GetAccountsDetail(v.input)
		assert.Equal(t, float32(v.accountDetails.AccountBalance), accountsDetails.AccountBalance)
	}
}
func TestCallGetAccountsDetailsForIdNotFound(t *testing.T) {
	accountsService := &AccountsService{}

	getAccountsDetail = func(accountId string) (*domain.AccountsSubsciptionDetail, *utils.ApplicationError) {
		return nil, &utils.ApplicationError{Message: "Account Id %v does not exists,please subscribe to use the service", StatusCode: http.StatusNotFound, Code: "accId_not_found"}
	}
	_, err := accountsService.GetAccountsDetail("prakhars")
	assert.Equal(t, 404, err.StatusCode)
}
func TestCallGetAccountsDetailsForDbError(t *testing.T) {
	accountsService := &AccountsService{}

	getAccountsDetail = func(accountId string) (*domain.AccountsSubsciptionDetail, *utils.ApplicationError) {
		return nil, &utils.ApplicationError{Message: fmt.Sprintf("Error in accessing db"), StatusCode: http.StatusInternalServerError, Code: "db_error"}
	}
	_, err := accountsService.GetAccountsDetail("prakhars")
	assert.Equal(t, "db_error", err.Code)
}

func TestCallGetPacksAndChannelsWithoutError(t *testing.T) {
	accountsService := &AccountsService{}

	getPacksAndChannels = func() (*domain.PacksAndChannels, *utils.ApplicationError) {
		packs := []domain.PacksDetails{{ID: 1, ChannelId: []int{10, 11, 12}, Name: "Silver pack", Price: 50}, {ID: 2, ChannelId: []int{10, 11, 12, 13, 14}, Name: "Gold pack", Price: 100}}
		channels := []domain.ChannelsDetails{{ID: 10, Name: "Zee", Price: 10}, {ID: 11, Name: "Sony", Price: 15}, {ID: 12, Name: "Star Plus", Price: 20}, {ID: 13, Name: "Discovery", Price: 10}, {ID: 14, Name: "NatGeo", Price: 20}}
		return &domain.PacksAndChannels{Packs: packs, Channels: channels}, nil
	}
	output, _ := accountsService.GetPacksAndChannels()
	assert.Equal(t, 2, len(output.Packs))
	assert.Equal(t, 5, len(output.Channels))
}

func TestCallGetPacksAndChannelsWithError(t *testing.T) {
	accountsService := &AccountsService{}

	getPacksAndChannels = func() (*domain.PacksAndChannels, *utils.ApplicationError) {
		return nil, &utils.ApplicationError{Message: fmt.Sprintf("Unable to get packs from DB"), StatusCode: http.StatusInternalServerError, Code: "db_error"}
	}
	_, err := accountsService.GetPacksAndChannels()
	assert.Equal(t, "db_error", err.Code)
}

func TestCallRechargeAmountsWithoutError(t *testing.T) {
	accountsService := &AccountsService{}
	rechargeAmount = func(amount int, accountId string) (float32, *utils.ApplicationError) {
		return 150, nil
	}
	output, _ := accountsService.RechargeAmount(150, "prakhars")
	assert.Equal(t, float32(150), output)
}

//Table driven test for RechargeAmount
var RechargeAmountData = []struct {
	amount       int
	accountId    string
	output       float32
	err          *utils.ApplicationError
	expectedCode string
}{
	{
		amount:       100,
		accountId:    "prakhars",
		output:       0,
		err:          &utils.ApplicationError{Message: "Account Id not subscribed", StatusCode: http.StatusNotFound, Code: "accId_not_found"},
		expectedCode: "accId_not_found",
	},
	{
		amount:       100,
		accountId:    "prakhars",
		output:       0,
		err:          &utils.ApplicationError{Message: "Error in accessing db", StatusCode: http.StatusInternalServerError, Code: "db_error"},
		expectedCode: "db_error",
	},
	{
		amount:       100,
		accountId:    "prakhars",
		output:       0,
		err:          &utils.ApplicationError{Message: "Error in updating db", StatusCode: http.StatusInternalServerError, Code: "db_error"},
		expectedCode: "db_error",
	},
}

func TestCallRechargeAmountsWithError(t *testing.T) {
	accountsService := &AccountsService{}
	for _, v := range RechargeAmountData {
		rechargeAmount = func(amount int, accountId string) (float32, *utils.ApplicationError) {
			return v.output, v.err
		}
		_, err := accountsService.RechargeAmount(v.amount, v.accountId)
		assert.Equal(t, v.expectedCode, err.Code)
	}
}

//Table driven test for UpdateEmailAndPhone
var UpdateEmailAndPhoneData = []struct {
	email        string
	phone        string
	accountId    string
	success      bool
	err          *utils.ApplicationError
	expectedCode bool
}{
	{
		email:        "prakhars@gmail.com",
		phone:        "1234567890",
		accountId:    "prakhars",
		success:      true,
		err:          nil,
		expectedCode: true,
	},
	{
		email:     "prakhars@gmail.com",
		phone:     "1234567890",
		accountId: "prakhars",
		success:   false,
		err: &utils.ApplicationError{
			Message:    "Error in updating db",
			StatusCode: http.StatusInternalServerError,
			Code:       "db_error",
		},
		expectedCode: false,
	},
	{
		email:     "prakhars@gmail.com",
		phone:     "1234567890",
		accountId: "prakhars",
		success:   false,
		err: &utils.ApplicationError{
			Message:    "Account Id not subscribed",
			StatusCode: http.StatusNotFound,
			Code:       "accId_not_found",
		},
		expectedCode: false,
	},
}

func TestCallUpdateEmailAndPhone(t *testing.T) {
	accountsService := &AccountsService{}

	for _, v := range UpdateEmailAndPhoneData {
		updateEmailAndPhone = func(string, string, string) (bool, *utils.ApplicationError) {
			return v.success, v.err
		}
		output, _ := accountsService.UpdateEmailAndPhone(v.email, v.phone, v.accountId)
		assert.Equal(t, v.expectedCode, output)
	}
}
