package entity

type UserId string

type User struct {
	Id          UserId
	Name        string
	PhoneNumber string
}
