package router

import (
	"examples/gin_xorm_viper/internal/interface/service"
	"examples/gin_xorm_viper/internal/pkg/utils"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/gin"
)

const IdAuthRouter = "router-auth"

type authRouter struct {
	gone.Flag
	gin.RouteGroup
	root  gin.RouteGroup     `gone:"*"`
	iUser service.IUserLogin `gone:"*"`
}

func (r *authRouter) GonerName() string {
	return IdAuthRouter
}

func (r *authRouter) Init() {
	r.RouteGroup = r.root.Group("/api", r.auth)
}

func (r *authRouter) auth(ctx *gin.Context, in struct {
	authorization string `gone:"http,header"`
}) error {
	token, err := utils.GetBearerToken(in.authorization)
	if err != nil {
		return gone.ToError(err)
	}
	userId, err := r.iUser.GetUserIdFromToken(token)
	utils.SetUserId(ctx, userId)
	return err
}
