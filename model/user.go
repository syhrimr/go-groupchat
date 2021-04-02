package model

import "time"

type User struct {
	UserID     int64     `json:"user_id"`
	Username   string    `json:"username"`
	Password   string    `json:"omitempty"`
	Salt       string    `json:"omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	ProfilePic string    `json:"profile_pic"`
}
