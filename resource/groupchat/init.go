package groupchat

import (
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/lolmourne/go-groupchat/model"
)

type RedisResource struct {
	rdb  *redis.Client
	next DBItf
}

type DBResource struct {
	db *sqlx.DB
}

type DBItf interface {
<<<<<<< HEAD
	GetJoinedRoom(userID int64) ([]model.Room,error)
	GetRoomByID(roomID int64) (model.Room,error)
	GetRooms(userID int64) ([]model.Room,error)
	CreateRoom(roomName string, adminID int64, description string, categoryID string) error
=======
	GetJoinedRoom(userID int64) ([]model.Room, error)
	GetRoomByID(roomID int64) (model.Room, error)
	GetRooms(userID int64) ([]model.Room, error)
	CreateRoom(roomName string, adminID int64, description string, categoryID int64) error
>>>>>>> ce4d1b61a1a0b0a0c256d047b7aaea7a6a1e004d
	AddRoomParticipant(roomID, userID int64) error
	LeaveRoom(roomID, userID int64) error
}

func NewRedisResource(rdb *redis.Client, next DBItf) DBItf {
	return &RedisResource{
		rdb:  rdb,
		next: next,
	}
}

func NewDBResource(dbParam *sqlx.DB) DBItf {
	return &DBResource{
		db: dbParam,
	}
}
