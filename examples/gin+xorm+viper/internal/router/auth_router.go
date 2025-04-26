package router

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/gin"
	"template_module/internal/interface/service"
	"template_module/internal/pkg/utils"
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
