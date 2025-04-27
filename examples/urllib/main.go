package main

import (
	"fmt"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/urllib"
	"github.com/imroc/req/v3"
)

type MyService struct {
	gone.Flag
	*req.Request `gone:"*"` // 注入 *req.Request
	//*req.Client `gone:"*"` // 注入 *req.Client

	//urllib.Client `gone:"*"` // 注入 urllib.Client 接口
}

func (s *MyService) GetData() (string, error) {
	// 发起 GET 请求
	resp, err := s.
		SetHeader("Accept", "application/json").
		Get("https://ipinfo.io")
	if err != nil {
		return "", err
	}

	// 获取响应内容
	return resp.String(), nil
}

func main() {
	gone.
		Load(&MyService{}).
		Loads(urllib.Load). // 加载 URLlib 组件
		Run(func(s *MyService) {
			data, err := s.GetData()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println("Data:", data)
		})
}
