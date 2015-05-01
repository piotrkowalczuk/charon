package model

import "time"

// User ...
type User struct {
	ID                int64
	Password          string
	Username          string
	FirstName         string
	LastName          string
	IsActive          bool
	IsStaff           bool
	IsSuperuser       bool
	IsConfirmed       bool
	ConfirmationToken string
	LastLoginAt       time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// NewUser creates user with default properties.
func NewUser(username, password, firstName, lastName, confirmationToken string) *User {
	now := time.Now()
	never := time.Unix(0, 0)

	return &User{
		Password:          password,
		Username:          username,
		FirstName:         firstName,
		LastName:          lastName,
		IsActive:          false,
		IsStaff:           false,
		IsSuperuser:       false,
		IsConfirmed:       false,
		ConfirmationToken: confirmationToken,
		CreatedAt:         now,
		LastLoginAt:       never,
		UpdatedAt:         never,
	}
}

// String ...
func (u *User) String() string {
	return u.FirstName + " " + u.LastName
}
