package router

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/gin"
)

const IdRouterPub = "router-pub"

type pubRouter struct {
	gone.Flag
	gin.IRouter
	root gin.RouteGroup `gone:"*"`
}

func (r *pubRouter) GonerName() string {
	return IdRouterPub
}

func (r *pubRouter) Init() {
	r.IRouter = r.root.Group("/api")
}
