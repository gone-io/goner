[//]: # (desc: simple web demo， 使用 gin, gorm, viper, mysql)
<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>

# Gone框架 Gin+GORM+Viper 集成示例

本示例展示了如何使用Gone框架与Gin、GORM和Viper组件集成，实现一个简单的Web应用程序。

## 项目概述

本示例演示了以下功能：

- 使用Gone框架的依赖注入机制
- 集成Gin框架实现Web路由和控制器
- 集成GORM框架实现数据库访问
- 使用Viper进行配置管理

## 项目结构

```
.
├── config/
│   └── default.properties  # 配置文件
├── go.mod                  # Go模块定义
└── main.go                 # 主程序
```

## 配置说明

配置文件`config/default.properties`包含MySQL数据库连接信息：

```properties
gorm.mysql.dsn=root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local
```

## 代码实现

### 主程序

`main.go`文件包含了整个应用的核心逻辑：

```go
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
```

## 代码解析

### 依赖注入

Gone框架使用标签`gone:"*"`进行依赖注入：

1. 在`HelloController`中注入了`gin.IRouter`和`UserRepository`
2. 在`UserRepository`中注入了`*gorm.DB`

### 控制器

`HelloController`实现了`gin.Controller`接口的`Mount`方法，注册了两个路由：

- `GET /hello`：返回一个简单的问候消息
- `GET /user/:id`：根据ID查询用户信息

### 数据模型和仓库

- `User`结构体定义了用户模型
- `UserRepository`提供了数据库访问方法，如`GetByID`

### 应用启动

在`main`函数中：

1. 使用`gone.Loads`加载基础组件：
   - `goner.BaseLoad`：基础组件
   - `goneGorm.Load`：GORM核心组件
   - `mysql.Load`：MySQL驱动
   - `gin.Load`：Gin组件

2. 加载自定义组件：
   - `&HelloController{}`：控制器
   - `&UserRepository{}`：仓库

3. 调用`Serve()`启动应用

## 运行示例

### 环境准备

1. 确保已安装Go环境（推荐Go 1.16+）
2. 准备MySQL数据库，创建`test`数据库
3. 根据需要修改`config/default.properties`中的数据库连接信息

### 创建数据表

在MySQL中执行以下SQL创建用户表：

```sql
CREATE TABLE `users` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 插入测试数据
INSERT INTO `users` (`id`, `name`) VALUES (1, 'Gone User');
```

### 启动应用

```bash
go run main.go
```

### 测试API

1. 访问问候接口：
   ```
   curl http://localhost:8080/hello
   ```
   预期返回：`"Hello, Gone!"`

2. 查询用户信息：
   ```
   curl http://localhost:8080/user/1
   ```
   预期返回：`{"ID":1,"Name":"Gone User"}`

## 总结

本示例展示了Gone框架如何与Gin、GORM和Viper集成，实现一个简单但功能完整的Web应用。通过Gone框架的依赖注入机制，各组件之间的耦合度降低，代码更加清晰和易于维护。

这个示例可以作为使用Gone框架开发Web应用的起点，您可以在此基础上扩展更多功能，如添加中间件、实现更复杂的业务逻辑等。