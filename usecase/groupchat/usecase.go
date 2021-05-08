package groupchat

import (
	"errors"

	"github.com/lolmourne/go-groupchat/model"
)

func (u UseCase) CreateGroupchat(name string, adminID int64, desc string, categoryID int64) (model.Room, error) {
	err := u.dbRoomRsc.CreateRoom(name, adminID, desc, categoryID)
	if err != nil {
		return model.Room{}, err
	}

	return model.Room{
		Name:        name,
		AdminUserID: adminID,
		Description: desc,
		CategoryID:  categoryID,
	}, nil
}

func (u UseCase) EditGroupchat(name, desc, categoryID string) (model.Room, error) {
	panic("TBC")
}

func (u UseCase) JoinRoom(roomID, userID int64) error {
	err := u.dbRoomRsc.AddRoomParticipant(roomID, userID)

	return err
}

func (u UseCase) GetRoomByID(roomID int64) (model.Room, error) {
	return u.dbRoomRsc.GetRoomByID(roomID)
}

func (u UseCase) GetRoomList(userID int64) ([]model.Room, error) {
	return u.dbRoomRsc.GetRooms(userID)
}

func (u UseCase) GetRoomByCategoryID(userID, categoryID int64) ([]model.Room, error) {
	return u.dbRoomRsc.GetRoomByCategoryID(userID, categoryID)
}

func (u UseCase) DeleteRoom(roomID, userID int64) error {
	room, err := u.dbRoomRsc.GetRoomByID(roomID)
	if err != nil {
		return err
	}

	if userID != room.AdminUserID {
		return errors.New("user is not an admin")
	}

	err = u.dbRoomRsc.DeleteRoom(roomID)
	if err != nil {
		return errors.New("failed to delete a room")
	}

	return nil
}
