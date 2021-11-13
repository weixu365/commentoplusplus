package model

import (
	"time"
)

type Email struct {
	Email                      string    `json:"email"`
	UnsubscribeSecretHex       string    `json:"unsubscribeSecretHex"`
	LastEmailNotificationDate  time.Time `json:"lastEmailNotificationDate"`
	SendReplyNotifications     bool      `json:"sendReplyNotifications"`
	SendModeratorNotifications bool      `json:"sendModeratorNotifications"`
}
