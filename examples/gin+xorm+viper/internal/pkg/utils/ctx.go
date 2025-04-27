package utils

import (
	"examples/gin_xorm_viper/internal/pkg/e"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/gin"
	"strings"
)

const UserIdKey = "user-id"

func SetUserId(ctx *gin.Context, userId int64) {
	ctx.Set(UserIdKey, userId)
}

func GetUserId(ctx *gin.Context) int64 {
	value, exists := ctx.Get(UserIdKey)

	if !exists {
		panic(gone.ToError("user id not found in context"))
	}
	return value.(int64)
}

func GetBearerToken(authorizationHeader string) (string, error) {
	if !strings.HasPrefix(authorizationHeader, "Bearer ") {
		return "", e.ErrUnauthorized
	}

	token := strings.TrimPrefix(authorizationHeader, "Bearer ")
	if token == "" {
		return "", e.ErrUnauthorized
	}
	return token, nil
}
