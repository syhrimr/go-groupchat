package model

import "time"

type Room struct {
	RoomID      int64     `json:"room_id"`
	Name        string    `json:"name"`
	AdminUserID int64     `json:"admin_user_id"`
	Description string    `json:"description"`
	CategoryID  int64     `json:"category_id"`
	CreatedAt   time.Time `json:"created_at"`
}
