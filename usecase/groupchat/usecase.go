package groupchat

import (
	"log"
	"strconv"

	"github.com/lolmourne/go-groupchat/model"
)

func (u UseCase) CreateGroupchat(name, adminID, desc, categoryID string) (model.Room, error) {
	//Room set to always return nil since unable to get lastInsertId
	admin, err := strconv.ParseInt(adminID, 10, 64)
	if err != nil {
		log.Println(err)
		return model.Room{}, err
	}

	room := u.dbRoomRsc.CreateRoom(name, admin, desc, categoryID)
	return model.Room{}, room

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

func (u UseCase) LeaveRoom(roomID, userID int64) error {
	return u.dbRoomRsc.LeaveRoom(roomID, userID)
}