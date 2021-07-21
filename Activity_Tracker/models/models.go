package models

type User struct{
	UserId string
	Name string 
	Email string
	Phone string
}

type Activity struct{
	UserId string
	ActivityType string 
	ActivityStatus string
	ActivityStartTime string
	ActivityEndTime string
	Duration string
	Valid bool
	ActivityLabel string
	ActivityId string
}