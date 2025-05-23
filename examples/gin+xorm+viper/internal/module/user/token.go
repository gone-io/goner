package user

import (
	"encoding/base64"
	"examples/gin_xorm_viper/internal/pkg/utils"
	"fmt"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/xorm"
	"time"
)

type iUserToken struct {
	gone.Flag

	db xorm.Engine `gone:"*"`
}

type TokenRecord struct {
	Id     int64 `xorm:"pk autoincr"`
	Token  string
	UserId uint64

	CreatedAt *time.Time `xorm:"created"`
	DeletedAt *time.Time `xorm:"deleted"`
}

func (s *iUserToken) CreateToken(userId uint64) (token string, err error) {
	key := fmt.Sprintf("%d:%s", userId, utils.CreateToken())
	token = base64.URLEncoding.EncodeToString([]byte(key))
	_, err = s.db.Insert(&TokenRecord{Token: token, UserId: userId})
	return token, gone.ToError(err)
}
func (s *iUserToken) ParseToken(token string) (userId uint64, err error) {
	var t TokenRecord
	has, err := s.db.Where("token = ?", token).Get(&t)
	if err != nil {
		return 0, gone.ToError(err)
	}
	if !has {
		return 0, gone.NewParameterError("token not found")
	}
	return t.UserId, nil
}

func (s *iUserToken) DestroyToken(token string) (err error) {
	_, err = s.db.Where("token = ?", token).Delete(&TokenRecord{})
	return gone.ToError(err)
}
