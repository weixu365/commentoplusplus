package model

import (
	"time"
)

type Comment struct {
	CommentHex   string    `json:"commentHex" db:"commentHex"`
	Domain       string    `json:"domain,omitempty" db:"domain"`
	Path         string    `json:"url,omitempty" db:"path"`
	CommenterHex string    `json:"commenterHex" db:"commenterHex"`
	Markdown     string    `json:"markdown" db:"markdown"`
	Html         string    `json:"html" db:"html"`
	ParentHex    string    `json:"parentHex" db:"parentHex"`
	Score        int       `json:"score" db:"score"`
	State        string    `json:"state,omitempty" db:"state"`
	CreationDate time.Time `json:"creationDate" db:"creationDate"`
	Direction    int       `json:"direction" db:"direction"`
	Deleted      bool      `json:"deleted" db:"deleted"`
}
