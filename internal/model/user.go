package model

import (
	"time"
)

type Role int32

type User struct {
	UUID      string
	Name      string
	Email     string
	Password  string
	Role      Role
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type CreateUser struct {
	Name            string
	Email           string
	Password        string
	PasswordConfirm string
	Role            Role
}

type UpdateUser struct {
	Name  *string
	Email *string
	Role  Role
}
