package svc

import (
	"ipashare/pkg/auth"
	"ipashare/pkg/e"
	"ipashare/pkg/tools"
)

type User struct {
	Service
}

func (u *User) Login(username, password string) (string, error) {
	user, err := u.store.User.Query(username)
	if err != nil {
		return "", e.NewWithStack(e.DBError, err)
	}
	if tools.MD5LowercaseEncode(password+user.Salt) != user.Password {
		return "", e.NewWithStack(e.ErrPassword, err)
	}
	jwt, err := auth.GenerateJWT(username)
	if err != nil {
		return "", e.NewWithStack(e.ErrTokenGen, err)
	}
	return jwt, nil
}

func (u *User) ChangePW(oldUsername, newUsername, password string) error {
	err := u.store.User.Update(oldUsername, newUsername, password)
	if err != nil {
		return e.NewWithStack(e.DBError, err)
	}
	auth.ChangeJwtSecret()
	return nil
}
