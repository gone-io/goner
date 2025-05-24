package service

type IUserToken interface {
	CreateToken(userId uint64) (token string, err error)
	ParseToken(token string) (userId uint64, err error)
	DestroyToken(token string) (err error)
}
