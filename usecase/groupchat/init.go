package groupchat

import (
	"github.com/lolmourne/go-groupchat/model"
	"github.com/lolmourne/go-groupchat/resource/groupchat"
)

type UseCase struct {
	dbRoomRsc groupchat.DBItf
}

type UsecaseItf interface {
	CreateGroupchat(name string, adminID int64, desc string, categoryID int64) (model.Room, error)
	EditGroupchat(name, desc, categoryID string) (model.Room, error)
	GetRoomByID(roomID int64) (model.Room, error)
	GetRoomList(userID int64) ([]model.Room, error)
	JoinRoom(roomID, userID int64) error
	LeaveRoom(roomID, userID int64) error
}

func NewUseCase(dbRsc groupchat.DBItf) UsecaseItf {
	return UseCase{
		dbRoomRsc: dbRsc,
	}

}
