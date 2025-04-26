package service

import "template_module/internal/interface/entity"

type IUserLogin interface {
	Register(registerParam *entity.RegisterParam) (*entity.LoginResult, error)

	Login(loginParam *entity.LoginParam) (*entity.LoginResult, error)
	Logout(token string) error

	GetUserIdFromToken(token string) (userId int64, err error)
}

type IUser interface {
	GetUserById(userId int64) (*entity.User, error)
}
