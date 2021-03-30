package groupchat

import "github.com/lolmourne/go-groupchat/model"

type UsecaseItf interface {
	CreateGroupchat(name, adminID, desc, categoryID string) (model.Room, error)
	EditGroupchat(name, desc, categoryID string) (model.Room, error)
	JoinRoom(roomID, userID int64) error
}
