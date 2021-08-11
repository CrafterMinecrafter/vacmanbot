package models

import "time"

type Penis struct {
	UserID  int       `json:"user_id"`
	Length  int       `json:"length"`
	Expires time.Time `json:"expires"`
}
