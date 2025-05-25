package service

import "examples/gin_xorm_viper/internal/interface/entity"

type IUserLogin interface {
	Register(registerParam *entity.RegisterParam) (*entity.LoginResult, error)

	Login(loginParam *entity.LoginParam) (*entity.LoginResult, error)
	Logout(token string) error

	GetUserIdFromToken(token string) (userId uint64, err error)
}

type IUser interface {
	GetUserById(userId uint64) (*entity.User, error)
}
