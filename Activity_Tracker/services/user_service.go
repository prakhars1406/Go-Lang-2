/*
This file is part of Services package.
This package deals with creating new user and getting user details.
This package also provides functions to fetch user activities
*/
package services

import (
	"log"
	"GitHub/Activity_Tracker/protoservices"
	"GitHub/Activity_Tracker/models"
	context "context"
	"time"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//This slice will store all the user details
var Users []models.User

//This slice will store all the users activities list
var Activity []models.Activity

type UserService struct {
}

/*
Returns NewUserResponse and error
Parameters Required
	- ctx : Context
	- in : Pointer to NewUserRequest data
Return Parameters
	- NewUserResponse : Pointer to NewUserResponse data
	- error : This returns error details in case of failure
Steps for function Flow
	1. This function will add new user to Users slice.
	2. Returns new user's details and user Id
Note:This function is exported function so can be used by other packages.
*/
func (s *UserService)AddUser(ctx context.Context, in *protoservices.NewUserRequest) (*protoservices.NewUserResponse, error){
	log.Printf("Welcome: %s", in.Name)
	userID := uuid.New()
	Users=append(Users,models.User{UserId:userID.String(),Name:in.Name,Email:in.Email,Phone:in.Phone})
	return &protoservices.NewUserResponse{UserId:userID.String(),Name:in.Name,Email:in.Email,Phone:in.Phone}, nil
}

/*
Returns QueryUserResponse and error
Parameters Required
	- ctx : Context
	- in : Pointer to QueryUserRequest data
Return Parameters
	- QueryUserResponse : Pointer to QueryUserResponse data
	- error : This returns error details in case of failure
Steps for function Flow
	1. This function will fetch user details based on user id
	2. Returns user's details and user Id
Note:This function is exported function so can be used by other packages.
*/
func (s *UserService)QueryUser(ctx context.Context, in *protoservices.QueryUserRequest) (*protoservices.QueryUserResponse, error){
	log.Printf("Welcome: %s", in.UserId)
	index:=-1
	for i,value:=range Users{
		if value.UserId==in.UserId{
			index=i
			break
		}
	}
	if index==-1{
		err := status.Error(codes.NotFound, "UserId  not found")
		return nil, err
	}
	return &protoservices.QueryUserResponse{UserId: Users[index].UserId,Name:Users[index].Name,Email:Users[index].Email,Phone:Users[index].Phone}, nil
}

/*
Returns QueryActivityResponse and error
Parameters Required
	- ctx : Context
	- in : Pointer to QueryActivityRequest data
Return Parameters
	- QueryActivityResponse : Pointer to QueryActivityResponse data
	- error : This returns error details in case of failure
Steps for function Flow
	1. This function will fetch user activity detail based on user id and activity id
	2. Returns user's activity details
Note:This function is exported function so can be used by other packages.
*/
func (s *UserService)QueryActivity(ctx context.Context, in *protoservices.QueryActivityRequest) (*protoservices.QueryActivityResponse, error){
	log.Printf("Welcome: %s", in.UserId)
	index:=-1
		for i,value:=range Activity{ 
			if value.UserId==in.GetUserId() && value.ActivityId==in.GetActivityId(){
				index=i
				break;
			}
		}
	if index==-1{
		err := status.Error(codes.NotFound, "UserId or Activity Id not found")
		return nil, err
	}
	return &protoservices.QueryActivityResponse{UserId:Activity[index].UserId,ActivityType:Activity[index].ActivityType,ActivityStatus:Activity[index].ActivityStatus,ActivityStartTime:Activity[index].ActivityStartTime,ActivityEndTime:Activity[index].ActivityEndTime,Duration:Activity[index].Duration,Valid:Activity[index].Valid,ActivityLabel:Activity[index].ActivityLabel,ActivityId:Activity[index].ActivityId}, nil
}


/*
Returns User's activity lists in stream format and error
Parameters Required
	- in : Pointer to QueryAllActivityRequest data
	- stream: stream to UserService_QueryAllActivityServer
Return Parameters
	- error : This returns error details in case of failure
Steps for function Flow
	1. This function will fetch user activity list in form of stream based on user id and activity id
	2. Returns user's activity details
Note:This function is exported function so can be used by other packages.
*/
func (s *UserService)QueryAllActivity(in *protoservices.QueryAllActivityRequest, stream protoservices.UserService_QueryAllActivityServer) error{
	log.Printf("Request from user:%s", in.UserId)
	for _,value:=range Activity{ 
		if value.UserId==in.GetUserId(){
			res := &protoservices.QueryAllActivityResponse{UserId:value.UserId,ActivityType:value.ActivityType,ActivityStatus:value.ActivityStatus,ActivityStartTime:value.ActivityStartTime,ActivityEndTime:value.ActivityEndTime,Duration:value.Duration,Valid:value.Valid,ActivityLabel:value.ActivityLabel,ActivityId:value.ActivityId}
			stream.Send(res)
			time.Sleep(1000 * time.Millisecond)
		}else{
			err := status.Error(codes.NotFound, "UserId not found")
			return err
		}
	}
	return nil
}