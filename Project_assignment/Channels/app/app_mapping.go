/*
This is the part of app package and this will be handle the user request and route it to suitable handler function.
*/
package app

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	domain "github.com/Channels/Domain"
	services "github.com/Channels/Services"
	utils "github.com/Channels/Utils"
)

type applicationMap struct{}

var (
	ApplicationMap applicationMapInterface
)

func init() {
	ApplicationMap = &applicationMap{}
}

type applicationMapInterface interface {
	getAccountDetails()
	callGetAccountsDetails(string) (*domain.AccountsSubsciptionDetail, *utils.ApplicationError)
	getPacksAndChannels()
	callGetPacksAndChannels() (*domain.PacksAndChannels, *utils.ApplicationError)
	rechargeAmount()
	callRechargeAmounts(int, string) (float32, *utils.ApplicationError)
	updateEmailAndPhone()
	callUpdateEmailAndPhone(string, string, string) (bool, *utils.ApplicationError)
	addChannels()
	callAddChannel(string, string) (float32, *utils.ApplicationError)
	subscribeToBase()
	callSubscribeToBase(string, string, int) (bool, bool, *utils.ApplicationError)
}

/*
Get the current balance and subciption details of the users
Parameters Required
	- No param

Return Parameters
	- No param

Steps for function Flow
	1. Takes user input
	2. Call account service function GetAccountsDetail
	3. Gets the response and check that function call was success or not
	4. Based on response show the appropriate console output

Note:This function will be only used by app functions thats why its not exported.
*/
func (appMap *applicationMap) getAccountDetails() {
	fmt.Println()
	fmt.Println("Please enter account Id")
	var accountId string
	fmt.Scanln(&accountId)
	if len(accountId) == 0 {
		log.Println("Account Id cannot be nil")
		return
	}
	accounts, err := ApplicationMap.callGetAccountsDetails(accountId)
	if err != nil {
		if err.Code == "accId_not_found" {
			fmt.Println("Account balance is 100 Rs.")
			fmt.Println("Current subscription: Not subscribed")
			return
		} else {
			fmt.Println(err.Message)
			return
		}
	}

	fmt.Printf("Account balance is %f Rs.\n", accounts.AccountBalance)
	fmt.Println("Current subscription:", accounts.CurrentSubscription)
	return
}
func (appMap *applicationMap) callGetAccountsDetails(accountId string) (*domain.AccountsSubsciptionDetail, *utils.ApplicationError) {
	accounts, err := services.AccountsService.GetAccountsDetail(accountId)
	return accounts, err
}

/*
Get the packs and channel details.
Parameters Required
	- No param

Return Parameters
	- No param

Steps for function Flow
	1. Call packs service function GetPacksAndChannels
	2. Gets the response and check that function call was success or not
	3. Based on response show the appropriate console output

Note:This function will be only used by app functions thats why its not exported.
*/
func (appMap *applicationMap) getPacksAndChannels() {

	packsAndChannels, err := ApplicationMap.callGetPacksAndChannels()
	if err != nil {
		log.Println(err.Message)
		return
	}
	channelMap := make(map[int]string)
	var channelBuilder strings.Builder
	for _, channels := range packsAndChannels.Channels {

		channelBuilder.WriteString(channels.Name + ": " + strconv.Itoa(channels.Price) + " Rs,")
		channelMap[channels.ID] = channels.Name
	}
	fmt.Println()
	fmt.Println("Available packs for subscription")
	for _, packs := range packsAndChannels.Packs {
		var b strings.Builder
		for _, presentChannels := range packs.ChannelId {
			b.WriteString(channelMap[presentChannels] + ",")
		}
		fmt.Println(packs.Name + ": " + b.String()[:len(b.String())-1] + ":" + strconv.Itoa(packs.Price) + " Rs")
	}
	fmt.Println("Available channels for subscription")
	fmt.Println(channelBuilder.String()[:len(channelBuilder.String())-1])
	return
}
func (appMap *applicationMap) callGetPacksAndChannels() (*domain.PacksAndChannels, *utils.ApplicationError) {
	packsAndChannels, err := services.PacksService.GetPacksAndChannels()
	return packsAndChannels, err
}

/*
Recharge user account balance.
Parameters Required
	- No param

Return Parameters
	- No param

Steps for function Flow
	1. Takes user input
	2. Call account service function RechargeAmount
	3. Gets the response and check that function call was success or not
	4. Based on response show the appropriate console output

Note:This function will be only used by app functions thats why its not exported.
*/
func (appMap *applicationMap) rechargeAmount() {
	fmt.Println()
	fmt.Println("Please enter recharge amount")
	var amount string
	fmt.Scanln(&amount)
	fmt.Println("Please enter account Id")
	var accountId string
	fmt.Scanln(&accountId)
	amt, err := strconv.Atoi(amount)
	if err != nil {
		fmt.Println("Please pass valid input")
		return
	}
	if amt == 0 {
		fmt.Println("Account Id cannot be nil")
		return
	}
	if len(accountId) == 0 {
		fmt.Println("account Id cannot be nil")
		return
	}
	balance, errr := ApplicationMap.callRechargeAmounts(amt, accountId)
	if errr != nil {
		fmt.Println(errr.Message)
		return
	}
	if balance != 0 {
		fmt.Println("Recharge completed successfully. Current balance is", balance)
		return
	}
	fmt.Println("Recharge Failed")
	return
}
func (appMap *applicationMap) callRechargeAmounts(amt int, accountId string) (float32, *utils.ApplicationError) {
	balance, err := services.AccountsService.RechargeAmount(amt, accountId)
	return balance, err
}

/*
Update users email and phone number to send him notifcations.
Parameters Required
	- No param

Return Parameters
	- No param

Steps for function Flow
	1. Takes user input
	2. Call notification service function UpdateEmailAndPhone
	3. Gets the response and check that function call was success or not
	4. Based on response show the appropriate console output

Note:This function will be only used by app functions thats why its not exported.
*/
func (appMap *applicationMap) updateEmailAndPhone() {
	fmt.Println()
	fmt.Println("Update email and phone number for notifications")
	fmt.Println("Enter your email")
	var email string
	fmt.Scanln(&email)
	fmt.Println("Enter your phone number")
	var phone string
	fmt.Scanln(&phone)
	fmt.Println("Enter your account Id")
	var accountId string
	fmt.Scanln(&accountId)
	if len(email) == 0 {
		fmt.Println("email cannot be nil")
		return
	}
	if len(phone) == 0 {
		fmt.Println("phone cannot be nil")
		return
	}
	if len(accountId) == 0 {
		fmt.Println("account Id cannot be nil")
		return
	}
	success, errr := ApplicationMap.callUpdateEmailAndPhone(email, phone, accountId)
	if errr != nil {
		log.Println(errr.Message)
		return
	}
	if success == true {
		fmt.Println("Email and Phone updated successfully.")
		return
	}
	fmt.Println("Email and Phone updation failed.")
	return
}
func (appMap *applicationMap) callUpdateEmailAndPhone(email string, phone string, accountId string) (bool, *utils.ApplicationError) {
	success, err := services.NotificationService.UpdateEmailAndPhone(email, phone, accountId)
	return success, err
}

/*
Add channels to users accounts.
Parameters Required
	- No param

Return Parameters
	- No param

Steps for function Flow
	1. Takes user input
	2. Call packs service function AddChannel
	3. Gets the response and check that function call was success or not
	4. Based on response show the appropriate console output

Note:This function will be only used by app functions thats why its not exported.
*/
func (appMap *applicationMap) addChannels() {
	fmt.Println()
	fmt.Println("Enter channel name to add")
	var channel string
	fmt.Scanln(&channel)
	fmt.Println("Enter your account Id")
	var accountId string
	fmt.Scanln(&accountId)
	if len(channel) == 0 {
		fmt.Println("channel cannot be nil")
		return
	}
	if len(accountId) == 0 {
		fmt.Println("accountId cannot be nil")
		return
	}
	balance, err := ApplicationMap.callAddChannel(channel, accountId)
	if err != nil {
		fmt.Println(err.Message)
		return
	}
	if balance != 0 {
		fmt.Println("Channels added successfully.")
		fmt.Printf("Account balance: %f Rs.\n", balance)
		return
	}
	fmt.Println("Adding channel failed")
	return
}
func (appMap *applicationMap) callAddChannel(channel string, accountId string) (float32, *utils.ApplicationError) {
	balance, err := services.PacksService.AddChannel(channel, accountId)
	return balance, err
}

/*
Subcribe new user to base pack.
Parameters Required
	- No param

Return Parameters
	- No param

Steps for function Flow
	1. Takes user input
	2. Call account service function SubscribeToBase
	3. Gets the response and check that function call was success or not
	4. Based on response show the appropriate console output

Note:This function will be only used by app functions thats why its not exported.
*/
func (appMap *applicationMap) subscribeToBase() {
	fmt.Println()
	fmt.Println("Enter your account Id")
	var accountId string
	fmt.Scanln(&accountId)
	fmt.Println("Enter the Pack you wish to subscribe: (Silver: S, Gold: G)")
	var pack string
	fmt.Scanln(&pack)
	fmt.Println("Enter the number of months for subsciption.")
	var months string
	fmt.Scanln(&months)
	if len(accountId) == 0 {
		fmt.Println("account Id cannot be nil")
		return
	}
	if strings.ToLower(pack) == "g" {

	} else if strings.ToLower(pack) == "s" {

	} else {
		fmt.Println("Invalid pack selected")
		return
	}
	month, err := strconv.Atoi(months)
	if err != nil {
		fmt.Println("Please pass valid months input")
		return
	}
	if month <= 0 {
		fmt.Println("Months should be greater then 0")
		return
	}
	success, notify, errr := ApplicationMap.callSubscribeToBase(accountId, pack, month)
	if errr != nil {
		fmt.Println(errr.Message)
		return
	}
	if success {
		printSubscribeMesg(pack, months, month, notify)
	}
	return
}
func (appMap *applicationMap) callSubscribeToBase(accountID string, pack string, months int) (bool, bool, *utils.ApplicationError) {
	success, notify, err := services.AccountsService.SubscribeToBase(accountID, pack, months)
	return success, notify, err
}

/*
Utility function used by subscribeToBase to print the formatted output
Parameters Required
	- pack : Pack which user has selected
	- months : Number of months in string format
	- month : Number of month in int format
	- notify : Takes bool input that user should get notification or not

Return Parameters
	- No param

Steps for function Flow
	1. Takes input params
	2. Checks subscribe pack is gold or silver
	3. Print console output based on subscribed pack

Note:This function will be only used by app functions thats why its not exported.
*/
func printSubscribeMesg(pack string, months string, month int, notify bool) {
	if strings.ToLower(pack) == "g" {
		fmt.Println()
		fmt.Println("You have successfully subscribed the following pack: Gold")
		fmt.Println("Monthly price: 100 Rs.")
		fmt.Println("No of months: " + months)
		totalPrice := month * 100
		fmt.Println("Subscription Amount: ", totalPrice)
		var discount float32
		if month >= 3 {
			discount = (float32(totalPrice) / 100) * 10
			fmt.Println("Discount applied: ", discount)
		}
		fmt.Println("Final Price after discount:", totalPrice-int(discount))
	} else if strings.ToLower(pack) == "s" {
		fmt.Println()
		fmt.Println("You have successfully subscribed the following pack: Silver")
		fmt.Println("Monthly price: 50 Rs.")
		totalPrice := month * 50
		fmt.Println("Subscription Amount: ", totalPrice, " Rs")
		fmt.Println("No of months: ", month)
		var discount float32
		if month >= 3 {
			discount = (float32(totalPrice) / 100) * 10
			fmt.Println("Discount applied: ", discount)
		}
		fmt.Println("Final Price after discount:", totalPrice-int(discount))
	}
	if notify == true {
		fmt.Println("Email notification sent successfully")
		fmt.Println("SMS notification sent successfully")
	} else {
		fmt.Println("Update your email and phone no to get notification")
	}
	fmt.Println()
}
