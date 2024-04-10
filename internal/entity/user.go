package entity

type UserId string

type User struct {
	Id          UserId
	Name        string
	PhoneNumber string
	Email       string
}

func NewUser(id UserId, name string, phoneNumber string, email string) User {
	return User{
		Id:          id,
		Name:        name,
		PhoneNumber: phoneNumber,
		Email:       email,
	}
}
