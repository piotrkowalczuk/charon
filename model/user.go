package model

import "time"

// User ...
type User struct {
	ID          int64
	Password    string
	Username    string
	FirstName   string
	LastName    string
	IsActive    bool
	IsStaff     bool
	IsSuperuser bool
	LastLoginAt *time.Time
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}
