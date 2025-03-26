package main

import (
	"context"
	"fmt"

	"github.com/gone-io/gone/v2"
	goneOpenAi "github.com/gone-io/goner/openai"
	"github.com/sashabaranov/go-openai"
)

// singleAIUser 演示单客户端使用
type singleAIUser struct {
	gone.Flag
	client *openai.Client `gone:"*"` // 默认客户端
}

func (s *singleAIUser) Use() {
	resp, err := s.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "你好！请用中文回答：什么是Gone框架？",
				},
			},
		},
	)
	if err != nil {
		fmt.Printf("创建聊天补全时出错: %v\n", err)
		return
	}

	fmt.Printf("AI回答: %s\n", resp.Choices[0].Message.Content)
}

// multiAIUser 演示多客户端使用
type multiAIUser struct {
	gone.Flag
	defaultClient *openai.Client `gone:"*"`        // 默认客户端
	baiduClient   *openai.Client `gone:"*,baidu"`  // 百度客户端
	aliyunClient  *openai.Client `gone:"*,aliyun"` // 阿里云客户端
}

func (m *multiAIUser) Use() {
	// 使用默认客户端
	resp, err := m.defaultClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "你好！请介绍一下你自己。",
				},
			},
		},
	)
	if err != nil {
		fmt.Printf("使用默认客户端时出错: %v\n", err)
		return
	}

	fmt.Printf("默认客户端回答: %s\n", resp.Choices[0].Message.Content)

	// 注意：这里仅作演示，实际使用时需要确保百度和阿里云客户端支持相同的API
	// 使用百度客户端
	resp, err = m.baiduClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "你好！请用中文介绍一下百度的AI服务。",
				},
			},
		},
	)
	if err != nil {
		fmt.Printf("使用百度客户端时出错: %v\n", err)
		return
	}

	fmt.Printf("百度客户端回答: %s\n", resp.Choices[0].Message.Content)
}

func main() {
	// 启动应用并加载OpenAI配置
	gone.
		NewApp(goneOpenAi.Load).
		Run(func(single *singleAIUser, multi *multiAIUser) {
			// 演示单客户端使用
			fmt.Println("=== 单客户端示例 ===")
			single.Use()

			fmt.Println("\n=== 多客户端示例 ===")
			multi.Use()
		})
}
