package model

import (
	"time"
)

type Commenter struct {
	CommenterHex string    `json:"commenterHex,omitempty"`
	Email        string    `json:"email,omitempty"`
	Name         string    `json:"name"`
	Link         string    `json:"link"`
	Photo        string    `json:"photo"`
	Provider     string    `json:"provider,omitempty"`
	JoinDate     time.Time `json:"joinDate,omitempty"`
	IsModerator  bool      `json:"isModerator"`
	Deleted      bool      `json:"deleted"`
}

type CommenterPassword struct {
	CommenterHex string    `json:"commenterHex,omitempty"`
	PasswordHash string    `json:"passwordHash,omitempty"`
	Email        string    `json:"email,omitempty"`
	Name         string    `json:"name"`
	Link         string    `json:"link"`
	Photo        string    `json:"photo"`
	Provider     string    `json:"provider,omitempty"`
	JoinDate     time.Time `json:"joinDate,omitempty"`
	IsModerator  bool      `json:"isModerator"`
	Deleted      bool      `json:"deleted"`
}
