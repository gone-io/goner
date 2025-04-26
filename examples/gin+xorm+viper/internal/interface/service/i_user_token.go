package service

type IUserToken interface {
	CreateToken(userId int64) (token string, err error)
	ParseToken(token string) (userId int64, err error)
	DestroyToken(token string) (err error)
}
