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
	GetJoinedRoom(userID int64) ([]model.Room, error)
	GetRoomByID(roomID int64) (model.Room, error)
	GetRooms(userID int64) ([]model.Room, error)
	CreateRoom(roomName string, adminID int64, description string, categoryID int64) error
	AddRoomParticipant(roomID, userID int64) error
	GetRoomByCategoryID(userID, categoryID int64) ([]model.Room, error)
	GetCategory() ([]model.Category, error)
	GetRoomParticipants(roomID int64) ([]model.User, error)
	LeaveRoom(userID, roomID int64) error
	DeleteRoom(roomID int64) error
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
