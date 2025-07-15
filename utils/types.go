package utils

type User struct {
	ID        string
	Provider  string
	FirstName string
	LastName  string
	Username  string
}

type userKey string

const UserKey userKey = "userKey"
