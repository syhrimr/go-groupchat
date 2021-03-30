package userauth

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"math/rand"

	"github.com/lolmourne/go-groupchat/model"
)

func (u *Usecase) Register(username, password, confirmPassword string) error {
	if confirmPassword != password {
		return errors.New("Confirm password is mismatched")
	}

	salt := RandStringBytes(32)
	password += salt

	h := sha256.New()
	h.Write([]byte(password))
	password = fmt.Sprintf("%x", h.Sum(nil))

	err := u.dbRsc.Register(username, password, salt)
	if err != nil {
		return err
	}

	return nil
}

func (u *Usecase) Login(username, password string) (*model.User, error) {
	user, err := u.dbRsc.GetUserByUserName(username)
	if err != nil {
		return nil, errors.New("user not found or password is incorrect")
	}

	password += user.Salt
	h := sha256.New()
	h.Write([]byte(password))
	hashedPassword := fmt.Sprintf("%x", h.Sum(nil))

	if user.Password != hashedPassword {
		return nil, errors.New("user not found or password is incorrect")
	}

	user.Password = ""
	user.Salt = ""

	return &user, nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
