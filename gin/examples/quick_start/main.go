package main

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/gin"
	"github.com/gone-io/goner/gin/injector"
	"reflect"
	"strings"
)

type HelloController struct {
	gone.Flag
	gin.IRouter `gone:"*"` // 注入路由器
}

// Mount 实现 gin.Controller 接口
func (h *HelloController) Mount() gin.MountError {
	h.GET("/hello", h.hello) // 注册路由
	h.GET("/token", func(token Token) string {
		return string(token)
	})

	type User struct {
		Id   uint64 `json:"id" form:"id"`
		Name string `json:"name" form:"name"`
	}

	h.GET("/user", func(query gin.Query[User]) User {
		return query.Get()
	})
	return nil
}

func (h *HelloController) hello() (string, error) {
	return "Hello, Gone!", nil
}

type Token string

func (t *Token) UserId() uint64 {
	return 1
}

var _ injector.TypeParser[*gin.Context] = (*tokenParser)(nil)

type tokenParser struct {
	gone.Flag
}

func (t tokenParser) Parse(context *gin.Context) (reflect.Value, error) {
	auth := context.GetHeader("Authorization")
	arr := strings.Split(auth, " ")
	if len(arr) == 2 && arr[0] == "Bearer" {
		token := Token(arr[1])
		return reflect.ValueOf(token), nil
	}
	return reflect.Value{}, gone.NewParameterError("invalid token")
}

func (t tokenParser) Type() reflect.Type {
	return reflect.TypeOf(Token(""))
}

func main() {
	gone.
		Load(&HelloController{}).
		Load(&tokenParser{}).
		Loads(gin.Load).
		Serve()
}

//curl -H "Authorization: Bearer 123456" http://127.0.0.1:8080/token
//curl http://127.0.0.1:8080/user?id=1&name=Gone
