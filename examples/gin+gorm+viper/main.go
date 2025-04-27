package main

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner"
	"github.com/gone-io/goner/gin"
	goneGorm "github.com/gone-io/goner/gorm"
	"github.com/gone-io/goner/gorm/mysql"
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

func main() {
	// 加载组件并启动应用
	gone.
		Loads(
			goner.BaseLoad,
			goneGorm.Load, // 加载 Gorm 核心组件
			mysql.Load,    // 加载 MySQL 驱动
			gin.Load,      // 加载 Gin 组件
		).
		Load(&HelloController{}). // 加载控制器
		Load(&UserRepository{}).  // 加载仓库
		Serve()
}
