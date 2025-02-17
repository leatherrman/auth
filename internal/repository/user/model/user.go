package model

import (
	"time"
)

type Role int32

type User struct {
	UUID           string     `db:"uuid"`
	Name           string     `db:"name"`
	Email          string     `db:"email"`
	Role           Role       `db:"role"`
	HashedPassword string     `db:"password"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      *time.Time `db:"updated_at"`
}
