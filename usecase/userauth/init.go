package userauth

import (
	"github.com/lolmourne/go-groupchat/resource"
)

type Usecase struct {
	dbRsc      resource.DBItf
	signingKey []byte
}

type UsecaseItf interface {
	Register(username, password, confirmPassword string) error
	Login(username, password string) (string, error)
	ValidateSession(accessToken string) (int64, error)
}

func NewUsecase(dbRsc resource.DBItf, signingKey string) UsecaseItf {
	return &Usecase{
		dbRsc:      dbRsc,
		signingKey: []byte(signingKey),
	}
}
