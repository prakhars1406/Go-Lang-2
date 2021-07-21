/*
This file is part of Services package.
This package deals with creating user activity such as play,read,eat and sleep.
This package also provides functions to validate the activity and check the status af activity
*/
package services

import (
	"log"
	"GitHub/Activity_Tracker/protoservices"
	"GitHub/Activity_Tracker/models"
	context "context"
	"github.com/google/uuid"
	"time"
	"strconv"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ActivityService struct {
}

/*
Returns valid status and duration
Parameters Required
	- startTime : Start time of an activity
	- endTime : End time of an activity
	- now : Current time in time.Time format
Return Parameters
	- bool : Validity status i.e true or false
	- string : Duration between Start time and end time in string format
Steps for function Flow
	1. This function will convert the start and end time in unix format
	2. Then it will find the difference.
	3. If difference is greater then 1 min then it will return valid=false else valid=true
	4. Returns valid status and duration
Note:This function is exported function so can be used by other packages.
*/
func Validate(startTime string,endTime string,now time.Time)(bool,string){
	log.Print(startTime)
	sTime, err := strconv.Atoi(startTime)
	if err != nil {
		sTime = 0
	}
	ts := int64(sTime)
	ts = ts * 1000000
	timeFromTS := time.Unix(0, ts)
	diff := now.Sub(timeFromTS)
	valid:=true
	if diff.Milliseconds() >int64(60000){
		valid=false
	}else{
		valid=true
	}
	duration:=diff.String()
	return valid,duration
}

/*
Returns PlayResponse and error
Parameters Required
	- ctx : Context
	- in : Pointer to PlayRequest data
Return Parameters
	- PlayResponse : Pointer to PlayResponse data
	- error : This returns error details in case of failure
Steps for function Flow
	1. This function will start play activity based on start flag.
	2. If user activity is not there and start=false then it will return error.
	3. If user activity is already there then it will close the user activity and update activity details
	4. Returns user's play details
Note:This function is exported function so can be used by other packages.
*/
func (s *ActivityService)Play(ctx context.Context, in *protoservices.PlayRequest) (*protoservices.PlayResponse, error){
	activityId := uuid.New().String()
	now := time.Now()
	unixNano := now.UnixNano()
	umillisec := unixNano / 1000000
	if in.GetStart(){
		Activity=append(Activity,models.Activity{UserId:in.GetUserId(),ActivityType:"Play",ActivityStatus:"ACTIVE",ActivityStartTime:strconv.FormatInt(umillisec, 10),ActivityEndTime:"",Duration:"",Valid:true,ActivityLabel:"play",ActivityId:activityId})
	}else{
		for i := len(Activity)-1; i >= 0; i-- {
			if Activity[i].ActivityType=="Play" && Activity[i].ActivityStatus=="ACTIVE"{
				log.Print("User Id found")
				startTime:=Activity[i].ActivityStartTime
				endTime:= strconv.FormatInt(umillisec, 10)
				log.Print(endTime)
				valid,duration:=Validate(startTime,endTime,now)
				Activity[i]=models.Activity{UserId:in.UserId,ActivityType:"Play",ActivityStatus:"DONE",ActivityStartTime:startTime,ActivityEndTime:endTime,Duration:duration,Valid:valid,ActivityLabel:"play",ActivityId:activityId}
				break;
			} else {
				err := status.Error(codes.NotFound, "No active Play  component found")
				return nil, err
			}
		}
	}
		index:=-1
		for i,value:=range Activity{ 
			if value.UserId==in.GetUserId() && value.ActivityId==activityId{
				index=i
				break;
			}
		}
		if index==-1{
			err := status.Error(codes.NotFound, "UserId or Activity Id not found")
			return nil, err
		}
	return &protoservices.PlayResponse{UserId:Activity[index].UserId,ActivityType:Activity[index].ActivityType,ActivityStatus:Activity[index].ActivityStatus,ActivityStartTime:Activity[index].ActivityStartTime,ActivityEndTime:Activity[index].ActivityEndTime,Duration:Activity[index].Duration,Valid:Activity[index].Valid,ActivityLabel:Activity[index].ActivityLabel,ActivityId:Activity[index].ActivityId}, nil

}

/*
Returns SleepResponse and error
Parameters Required
	- ctx : Context
	- in : Pointer to SleepRequest data
Return Parameters
	- SleepResponse : Pointer to PlayResponse data
	- error : This returns error details in case of failure
Steps for function Flow
	1. This function will start sleep activity based on start flag.
	2. If user activity is not there and start=false then it will return error.
	3. If user activity is already there then it will close the user activity and update activity details
	2. Returns user's sleep details
Note:This function is exported function so can be used by other packages.
*/
func (s *ActivityService)Sleep(ctx context.Context, in *protoservices.SleepRequest) (*protoservices.SleepResponse, error){
	activityId := uuid.New().String()
	now := time.Now()
	unixNano := now.UnixNano()
	umillisec := unixNano / 1000000
	if in.GetStart(){
		Activity=append(Activity,models.Activity{UserId:in.GetUserId(),ActivityType:"Sleep",ActivityStatus:"ACTIVE",ActivityStartTime:strconv.FormatInt(umillisec, 10),ActivityEndTime:"",Duration:"",Valid:true,ActivityLabel:"sleep",ActivityId:activityId})
	}else{
		for i := len(Activity)-1; i >= 0; i-- {
			if Activity[i].ActivityType=="Sleep" && Activity[i].ActivityStatus=="ACTIVE"{
				log.Print("User Id found")
				startTime:=Activity[i].ActivityStartTime
				endTime:= strconv.FormatInt(umillisec, 10)
				log.Print(endTime)
				valid,duration:=Validate(startTime,endTime,now)
				Activity[i]=models.Activity{UserId:in.UserId,ActivityType:"Sleep",ActivityStatus:"DONE",ActivityStartTime:startTime,ActivityEndTime:endTime,Duration:duration,Valid:valid,ActivityLabel:"sleep",ActivityId:activityId}
				break;
			} else {
				err := status.Error(codes.NotFound, "No active Sleep  component found")
				return nil, err
			}
		}
	}
		index:=-1
		for i,value:=range Activity{ 
			if value.UserId==in.GetUserId() && value.ActivityId==activityId{
				index=i
				break;
			}
		}
		if index==-1{
			err := status.Error(codes.NotFound, "UserId or Activity Id not found")
			return nil, err
		}
	return &protoservices.SleepResponse{UserId:Activity[index].UserId,ActivityType:Activity[index].ActivityType,ActivityStatus:Activity[index].ActivityStatus,ActivityStartTime:Activity[index].ActivityStartTime,ActivityEndTime:Activity[index].ActivityEndTime,Duration:Activity[index].Duration,Valid:Activity[index].Valid,ActivityLabel:Activity[index].ActivityLabel,ActivityId:Activity[index].ActivityId}, nil

}

/*
Returns EatResponse and error
Parameters Required
	- ctx : Context
	- in : Pointer to EatRequest data
Return Parameters
	- EatResponse : Pointer to EatResponse data
	- error : This returns error details in case of failure
Steps for function Flow
	1. This function will start eat activity based on start flag.
	2. If user activity is not there and start=false then it will return error.
	3. If user activity is already there then it will close the user activity and update activity details
	4. Returns user's eat details
Note:This function is exported function so can be used by other packages.
*/
func (s *ActivityService)Eat(ctx context.Context, in *protoservices.EatRequest) (*protoservices.EatResponse, error){
	activityId := uuid.New().String()
	now := time.Now()
	unixNano := now.UnixNano()
	umillisec := unixNano / 1000000
	if in.GetStart(){
		Activity=append(Activity,models.Activity{UserId:in.GetUserId(),ActivityType:"Eat",ActivityStatus:"ACTIVE",ActivityStartTime:strconv.FormatInt(umillisec, 10),ActivityEndTime:"",Duration:"",Valid:true,ActivityLabel:"eat",ActivityId:activityId})
	}else{
		for i := len(Activity)-1; i >= 0; i-- {
			if Activity[i].ActivityType=="Eat" && Activity[i].ActivityStatus=="ACTIVE"{
				log.Print("User Id found")
				startTime:=Activity[i].ActivityStartTime
				endTime:= strconv.FormatInt(umillisec, 10)
				log.Print(endTime)
				valid,duration:=Validate(startTime,endTime,now)
				Activity[i]=models.Activity{UserId:in.UserId,ActivityType:"Eat",ActivityStatus:"DONE",ActivityStartTime:startTime,ActivityEndTime:endTime,Duration:duration,Valid:valid,ActivityLabel:"eat",ActivityId:activityId}
				break;
			} else {
				err := status.Error(codes.NotFound, "No active Eat component found")
				return nil, err
			}
		}
	}
		index:=-1
		for i,value:=range Activity{ 
			if value.UserId==in.GetUserId() && value.ActivityId==activityId{
				index=i
				break;
			}
		}
		if index==-1{
			err := status.Error(codes.NotFound, "UserId or Activity Id not found")
			return nil, err
		}
	return &protoservices.EatResponse{UserId:Activity[index].UserId,ActivityType:Activity[index].ActivityType,ActivityStatus:Activity[index].ActivityStatus,ActivityStartTime:Activity[index].ActivityStartTime,ActivityEndTime:Activity[index].ActivityEndTime,Duration:Activity[index].Duration,Valid:Activity[index].Valid,ActivityLabel:Activity[index].ActivityLabel,ActivityId:Activity[index].ActivityId}, nil

}

/*
Returns ReadResponse and error
Parameters Required
	- ctx : Context
	- in : Pointer to ReadRequest data
Return Parameters
	- ReadResponse : Pointer to ReadResponse data
	- error : This returns error details in case of failure
Steps for function Flow
	1. This function will start read activity based on start flag.
	2. If user activity is not there and start=false then it will return error.
	3. If user activity is already there then it will close the user activity and update activity details
	4. Returns user's read details
Note:This function is exported function so can be used by other packages.
*/
func (s *ActivityService)Read(ctx context.Context, in *protoservices.ReadRequest) (*protoservices.ReadResponse, error){
	activityId := uuid.New().String()
	now := time.Now()
	unixNano := now.UnixNano()
	umillisec := unixNano / 1000000
	if in.GetStart(){
		Activity=append(Activity,models.Activity{UserId:in.GetUserId(),ActivityType:"Read",ActivityStatus:"ACTIVE",ActivityStartTime:strconv.FormatInt(umillisec, 10),ActivityEndTime:"",Duration:"",Valid:true,ActivityLabel:"read",ActivityId:activityId})
	}else{
		for i := len(Activity)-1; i >= 0; i-- {
			if Activity[i].ActivityType=="Read" && Activity[i].ActivityStatus=="ACTIVE"{
				log.Print("User Id found")
				startTime:=Activity[i].ActivityStartTime
				endTime:= strconv.FormatInt(umillisec, 10)
				log.Print(endTime)
				valid,duration:=Validate(startTime,endTime,now)
				Activity[i]=models.Activity{UserId:in.UserId,ActivityType:"Read",ActivityStatus:"DONE",ActivityStartTime:startTime,ActivityEndTime:endTime,Duration:duration,Valid:valid,ActivityLabel:"read",ActivityId:activityId}
				break;
			} else {
				err := status.Error(codes.NotFound, "No active Read component found")
				return nil, err
			}
		}
	}
		index:=-1
		for i,value:=range Activity{ 
			if value.UserId==in.GetUserId() && value.ActivityId==activityId{
				index=i
				break;
			}
		}
		if index==-1{
			err := status.Error(codes.NotFound, "UserId or Activity Id not found")
			return nil, err
		}
	return &protoservices.ReadResponse{UserId:Activity[index].UserId,ActivityType:Activity[index].ActivityType,ActivityStatus:Activity[index].ActivityStatus,ActivityStartTime:Activity[index].ActivityStartTime,ActivityEndTime:Activity[index].ActivityEndTime,Duration:Activity[index].Duration,Valid:Activity[index].Valid,ActivityLabel:Activity[index].ActivityLabel,ActivityId:Activity[index].ActivityId}, nil

}

/*
Returns IsValidResponse and error
Parameters Required
	- ctx : Context
	- in : Pointer to IsValidRequest data
Return Parameters
	- IsValidResponse : Pointer to ReadResponse data
	- error : This returns error details in case of failure
Steps for function Flow
	1. This function will check valid status of user's activity based on user Id and activity Id.
	2. If user activity is not found then it will return error.
	3. Returns user's activity valid status
Note:This function is exported function so can be used by other packages.
*/
func (s *ActivityService)IsValid(ctx context.Context, in *protoservices.IsValidRequest) (*protoservices.IsValidResponse, error){
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
		return &protoservices.IsValidResponse{Status:Activity[index].Valid},nil
}

/*
Returns IsDoneResponse and error
Parameters Required
	- ctx : Context
	- in : Pointer to IsDoneRequest data
Return Parameters
	- IsDoneResponse : Pointer to IsDoneResponse data
	- error : This returns error details in case of failure
Steps for function Flow
	1. This function will check user's activity completion status based on user Id and activity Id.
	2. If user activity is not found then it will return error.
	3. Returns user's activity completion status
Note:This function is exported function so can be used by other packages.
*/
func (s *ActivityService)IsDone(ctx context.Context, in *protoservices.IsDoneRequest) (*protoservices.IsDoneResponse, error){
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
		return &protoservices.IsDoneResponse{Status:Activity[index].ActivityStatus},nil
}