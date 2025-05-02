package main

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/gin"
	"gorm.io/gorm"
)

// 定义控制器
type HelloController struct {
	gone.Flag
	gin.IRouter `gone:"*"`      // 注入路由器
	uR          *UserRepository `gone:"*"`
}

// Mount 实现 gin.Controller 接口
func (h *HelloController) Mount() gin.MountError {
	h.GET("/hello", h.hello) // 注册路由
	h.GET("/user/:id", h.getUser)
	return nil
}

func (h *HelloController) hello() (string, error) {
	return "Hello, Gone!", nil
}
func (h *HelloController) getUser(in struct {
	id uint `param:"id"`
}) (*User, error) {

	user, err := h.uR.GetByID(in.id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// 定义数据模型和仓库
type User struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

type UserRepository struct {
	gone.Flag
	*gorm.DB `gone:"*"`
}

func (r *UserRepository) GetByID(id uint) (*User, error) {
	var user User
	err := r.First(&user, id).Error
	return &user, err
}

func LoadControllerAndRepository(loader gone.Loader) error {
	loader.
		MustLoad(&HelloController{}).
		MustLoad(&UserRepository{})
	return nil
}

func main() {
	// 加载组件并启动应用
	gone.
		Serve()
}
