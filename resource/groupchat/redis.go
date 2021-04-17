package groupchat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/lolmourne/go-groupchat/model"
)

func (dbr *RedisResource) GetJoinedRoom(userID int64) ([]model.Room, error) {
	val, err := dbr.rdb.Get(context.Background(), fmt.Sprintf("roomJoined:%d", userID)).Result()
	if err != nil {
		rooms, err := dbr.next.GetJoinedRoom(userID)
		if err != nil {
			return rooms, errors.New("no rooms found")
		}
		roomJSON, err := json.Marshal(rooms)
		if err != nil {
			return rooms, errors.New("fail to marshall")
		}
		stats := dbr.rdb.Set(context.Background(), fmt.Sprintf("roomJoined:%d", userID), roomJSON, time.Duration(0))
		log.Println(stats.Result())

		return rooms, err
	}

	//return dbr.next.GetJoinedRoom(userID)

	var rooms []model.Room
	json.Unmarshal([]byte(val), &rooms)

	return rooms, nil
}

func (dbr *RedisResource) CreateRoom(roomName string, adminID int64, description string, categoryID string) error {
	//no need redis for DML (insert)
	return dbr.next.CreateRoom(roomName, adminID, description, categoryID)
}

func (dbr *RedisResource) AddRoomParticipant(roomID, userID int64) error {
	//no need redis since DML query (insert)
	err := dbr.next.AddRoomParticipant(roomID, userID)
	if err == nil {
		resp := dbr.rdb.Del(context.Background(), fmt.Sprintf("roomJoined:%d", userID))
		if resp.Err() != nil {
			return err
		}
	}

	return err
}

func (dbr *RedisResource) GetRoomByID(roomID int64) (model.Room, error) {
	return dbr.next.GetRoomByID(roomID)
}
