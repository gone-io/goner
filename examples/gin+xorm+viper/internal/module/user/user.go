package user

import (
	"github.com/gone-io/goner/xorm"
	"template_module/internal/interface/entity"
	"template_module/internal/interface/service"
	"template_module/internal/pkg/e"
	"template_module/internal/pkg/utils"

	"github.com/gone-io/gone/v2"
)

type iUser struct {
	gone.Flag
	db  xorm.Engine `gone:"*"`
	log gone.Logger `gone:"*"`

	iDep       service.IDependent `gone:"*"`
	iUserToken service.IUserToken `gone:"*"`
}

func (s *iUser) Register(registerParam *entity.RegisterParam) (*entity.LoginResult, error) {
	err := s.iDep.DoSomething()
	if err != nil {
		return nil, gone.ToError(err)
	}

	exist, err := s.db.Where("username = ?", registerParam.Username).Exist(&entity.User{})
	if err != nil {
		return nil, gone.ToError(err)
	}
	if exist {
		return nil, e.UsernameHasBenUsed
	}

	passwordHash, err := utils.GenPasswordHash(registerParam.Password)
	if err != nil {
		return nil, gone.ToError(err)
	}

	user := entity.User{
		Username: registerParam.Username,
		Password: passwordHash,
	}

	_, err = s.db.Insert(&user)
	if err != nil {
		return nil, gone.ToError(err)
	}

	token, err := s.iUserToken.CreateToken(user.Id)
	if err != nil {
		return nil, gone.ToError(err)
	}

	return &entity.LoginResult{
		Token: token,
		User:  &user,
	}, nil
}

func (s *iUser) GetUserIdFromToken(token string) (userId int64, err error) {
	userId, err = s.iUserToken.ParseToken(token)
	if err != nil {
		s.log.Warnf("parse token failed, token: %s, err: %v", token, err)
		return 0, e.ErrUnauthorized
	}
	return
}

func (s *iUser) Login(loginParam *entity.LoginParam) (*entity.LoginResult, error) {
	var user entity.User
	has, err := s.db.Where("username = ?", loginParam.Username).Get(&user)
	if err != nil {
		return nil, gone.ToError(err)
	}
	if !has {
		return nil, e.LoginError
	}

	if !utils.VerifyPassword(user.Password, loginParam.Password) {
		return nil, e.LoginError
	}

	token, err := s.iUserToken.CreateToken(user.Id)
	if err != nil {
		return nil, gone.ToError(err)
	}

	return &entity.LoginResult{
		Token: token,
		User:  &user,
	}, nil
}

func (s *iUser) Logout(token string) error {
	return s.iUserToken.DestroyToken(token)
}

func (s *iUser) GetUserById(userId int64) (*entity.User, error) {
	var user entity.User
	has, err := s.db.ID(userId).Get(&user)
	if err != nil {
		return nil, gone.ToError(err)
	}
	if !has {
		return nil, e.UserNotFound
	}

	return &user, nil
}
