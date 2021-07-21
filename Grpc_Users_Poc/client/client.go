package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"google.golang.org/grpc"

	"GitHub/Grpc_Users_Poc/client/protoservices"
	"GitHub/Grpc_Users_Poc/client/services"
)

func main() {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	registerServiceClient := protoservices.NewRegisterServiceClient(conn)
	getUserServiceClient := protoservices.NewGetUserServiceClient(conn)
	getAavailableServiceClient := protoservices.NewGetAavailableServiceClient(conn)
	addServiceClient := protoservices.NewAddServiceClient(conn)
	checkUserServiceClient := protoservices.NewCheckUserServiceClient(conn)

	message := `Welcome to IT service
	1. Register
	2. Login
	3. Get Services
	4. Add Services
	5. Check subscibed Service.
	6. Exit`
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
			services.Register(registerServiceClient)
		case 2:
			services.GetUser(getUserServiceClient)
		case 3:
			services.GetService(getAavailableServiceClient)
		case 4:
			services.AddService(addServiceClient)
		case 5:
			services.CheckUserService(checkUserServiceClient)
		case 6:
			fmt.Println("Thanks for using our service")
			os.Exit(1)
		default:
			fmt.Println("Wrong option selected")
		}
	}

}
