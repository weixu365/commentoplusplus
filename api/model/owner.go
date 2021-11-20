package model

import (
	"time"
)

type Owner struct {
	OwnerHex       string    `json:"ownerHex"`
	Email          string    `json:"email"`
	PasswordHash   string    `json:"-"`
	Name           string    `json:"name"`
	ConfirmedEmail bool      `json:"confirmedEmail"`
	JoinDate       time.Time `json:"joinDate"`
}
