package router

import (
	"examples/gin_xorm_viper/internal/interface/entity"
	"examples/gin_xorm_viper/internal/interface/service"
	"examples/gin_xorm_viper/internal/pkg/utils"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/gin"
	"github.com/gone-io/goner/gin/injector"
	"reflect"
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

var _ injector.TypeParser[*gin.Context] = (*tokenParser)(nil)

type tokenParser struct {
	gone.Flag
}

func (t tokenParser) Parse(context *gin.Context) (reflect.Value, error) {
	userId := utils.GetUserId(context)
	return reflect.ValueOf(entity.Token{UserId: userId}), nil
}

func (t tokenParser) Type() reflect.Type {
	return reflect.TypeOf(entity.Token{})
}
