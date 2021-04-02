package groupchat

import (
	"database/sql"
	"time"
)

type RoomDB struct {
	RoomID      sql.NullInt64  `db:"room_id""`
	AdminUserID sql.NullInt64  `db:"admin_user_id"`
	Description sql.NullString `db:"description"`
	CategoryID  sql.NullInt64  `db:"category_id"`
	CreatedAt   time.Time      `db:"created_at"`
}
