/*
This is the app package and this will be handle the user request and route it to suitable handler function.
*/
package app

import (
	"fmt"
	"os"
	"strconv"
)

type applicationStart struct{}

var (
	ApplicationStart applicationStart
)

/*
This will show all the operation user can perform and based on user input it will call the handler function of mapping.
Parameters Required
	-No Param

Return Parameters
	- No Param

Steps for function Flow
	1. Show all options.
	2. Takes user input.
	3. Based on user input perform operation.

Note:This function is exported function so can be used by other packages.
*/
func (appStart *applicationStart) StartApp() {
	message := `Welcome to SatTV
	1. View account balance and current subscription
	2. Recharge Account
	3. View available packs and channels
	4. Subscribe to base packs
	5. Add channels to an existing subscription
	6. Update email and phone number for notifications
	7. Exit`
	for {
		fmt.Println()
		fmt.Println(message)
		var text string
		fmt.Scanln(&text)
		userInput, err := strconv.Atoi(text)
		if err != nil {
			fmt.Println("Please enter only number input")
			continue
		}
		switch userInput {
		case 1:
			ApplicationMap.getAccountDetails()
		case 2:
			ApplicationMap.rechargeAmount()
		case 3:
			ApplicationMap.getPacksAndChannels()
		case 4:
			ApplicationMap.subscribeToBase()
		case 5:
			ApplicationMap.addChannels()
		case 6:
			ApplicationMap.updateEmailAndPhone()
		case 7:
			fmt.Println("Thanks for using SatTV")
			os.Exit(1)
		default:
			fmt.Println("Wrong option selected")
		}
	}
}
