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

type Room struct {
	RoomID      int64     `json:"room_id"`
	AdminUserID int64     `json:"admin_user_id"`
	Description string    `json:"description"`
	CategoryID  int64     `json:"category_id"`
	CreatedAt   time.Time `json:"created_at"`
}
