package userauth

import (
	"time"

	"github.com/lolmourne/go-groupchat/model"
)

type Usecase struct {
	endpoint string
	timeout  time.Duration
}

type UsecaseItf interface {
	GetUserInfo(accessToken string) *model.User
}

func NewUsecase(endpoint string, timeout time.Duration) UsecaseItf {
	return &Usecase{
		endpoint: endpoint,
		timeout:  timeout,
	}
}
