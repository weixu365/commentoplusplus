package model

import "time"

type Moderator struct {
	Email   string    `json:"email"`
	Domain  string    `json:"domain"`
	AddDate time.Time `json:"addDate"`
}
