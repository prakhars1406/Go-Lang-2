package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"GitHub/Grpc_Users_Poc/client/protoservices"

	"google.golang.org/grpc/metadata"
)

func CheckUserService(c protoservices.CheckUserServiceClient) {
	fmt.Println()
	md := metadata.Pairs("token", "valid-token")
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	stream, err := c.CheckUserService(ctx)

	if err != nil {
		log.Fatalf("Error while opening stream and calling CheckServices: %v", err)
	}

	waitc := make(chan struct{})

	// send go routine
	go func() {
		for i := 0; i < 5; i++ {
			result := "Service_" + strconv.Itoa(i)
			stream.Send(&protoservices.CheckUserServiceRequest{
				ServiceName: result,
			})
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	// receive go routine
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Problem while reading server stream: %v", err)
				break
			}
			message := res.GetMessage()
			fmt.Printf("Received a new message of...: %s\n", message)
		}
		close(waitc)
	}()
	<-waitc
}
