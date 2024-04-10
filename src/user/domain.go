package user

import "github.com/google/uuid"

type UserDomain struct {
	Id        uuid.UUID
	Email     string
	FirstName string
	LastName  string
}

type UserPlatformDomain struct {
	Id           uuid.UUID
	PlatformName string
	PlatformID   string
	UserID       uuid.UUID
}
