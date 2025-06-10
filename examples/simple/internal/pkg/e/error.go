package e

import (
	"net/http"

	"github.com/gone-io/gone/v2"
)

var (
	ErrUnauthorized = gone.NewError(1001, "没有登录或登录失效", http.StatusUnauthorized)

	UsernameHasBenUsed = gone.NewParameterError("用户名已被使用", 1002)

	LoginError = gone.NewParameterError("用户名或密码错误", 1003)

	UserNotFound = gone.NewParameterError("用户不存在", 1004)
)
