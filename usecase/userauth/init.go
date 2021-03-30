package userauth

import (
	"github.com/lolmourne/go-groupchat/model"
	"github.com/lolmourne/go-groupchat/resource"
)

type Usecase struct {
	dbRsc resource.DBItf
}

type UsecaseItf interface {
	Register(username, password, confirmPassword string) error
	Login(username, password string) (*model.User, error)
}

func NewUsecase(dbRsc resource.DBItf) UsecaseItf {
	return &Usecase{
		dbRsc: dbRsc,
	}
}
