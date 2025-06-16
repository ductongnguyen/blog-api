package models

import (
	"fmt"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required"`
	Role      UserRole  `json:"role" validate:"required" gorm:"type:varchar(50)"` // Use the UserRole type, specify DB column type
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)

func (r UserRole) String() string {
	return string(r)
}

var validUserRoles = map[UserRole]struct{}{
	RoleAdmin: {},
	RoleUser:  {},
}

func (r UserRole) IsValid() bool {
	_, ok := validUserRoles[r]
	return ok
}

func ParseUserRole(s string) (UserRole, error) {
	role := UserRole(s)
	if !role.IsValid() {
		return "", fmt.Errorf("invalid user role: %s", s)
	}
	return role, nil
}
