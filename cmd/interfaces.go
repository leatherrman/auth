package main

import "time"

// Role is ...
type Role uint8

// NewUserData is ...
type NewUserData struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Role            uint8  `json:"role"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

// UserData is ...
type UserData struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Role      Role       `json:"role"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// UpdateUserData is ...
type UpdateUserData struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  Role   `json:"role"`
}

// ResponseUserID is ...
type ResponseUserID struct {
	ID int `json:"id"`
}
