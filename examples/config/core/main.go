package main

import (
	"fmt"
	"github.com/gone-io/gone/v2"
	"os"
)

// 定义一个复杂类型的配置项
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// 主配置结构体
type AppConfig struct {
	gone.Flag // 这是Gone框架的标识，表示这个结构体参与依赖注入

	// 应用名称配置：键名为app.name，默认值为my-app
	AppName string `gone:"config,app.name=my-app"`

	// 端口配置：键名为app.port，默认值为8080
	Port int `gone:"config,app.port=8080"`

	// 环境配置：键名为app.env，无默认值（如果未设置环境变量，将为空字符串）
	Environment string `gone:"config,app.env"`

	// 数据库配置：键名为app.database，这是一个复杂对象
	Database *DatabaseConfig `gone:"config,app.database"`
}

func main() {
	//设置环境变量
	_ = os.Setenv("GONE_APP_NAME", "my-awesome-app")
	_ = os.Setenv("GONE_APP_PORT", "9000")
	_ = os.Setenv("GONE_APP_ENV", "production")
	_ = os.Setenv("GONE_APP_DATABASE", `{"host":"localhost","port":5432,"username":"admin","password":"secret123"}`)

	var config AppConfig

	gone.
		Load(&config). // 加载配置结构体
		Run(func() {   // 运行应用逻辑
			fmt.Printf("应用名称: %s\n", config.AppName)
			fmt.Printf("运行端口: %d\n", config.Port)
			fmt.Printf("运行环境: %s\n", config.Environment)

			if config.Database != nil {
				fmt.Printf("数据库主机: %s:%d\n", config.Database.Host, config.Database.Port)
			}
		})
}
