package main

import (
	"context"
	"github.com/gone-io/gone/v2"
	"go.opentelemetry.io/otel/trace"
)

type YourComponent struct {
	gone.Flag
	tracer trace.Tracer `gone:"*,otel-tracer"` // 注入 OpenTelemetry Tracer
}

func (c *YourComponent) HandleRequest(ctx context.Context) {
	// tracer := otel.Tracer("otel-tracer")
	tracer := c.tracer

	// 创建新的 Span
	ctx, span := tracer.Start(ctx, "handle-request")
	// 确保在函数结束时结束 Span
	defer span.End()

	// 记录事件
	span.AddEvent("开始处理请求")

	// 处理业务逻辑...

	// 记录错误（如果有）
	// span.RecordError(err)
	// span.SetStatus(codes.Error, "处理请求失败")

	// 正常情况下设置状态为成功
	// span.SetStatus(codes.Ok, "")
}
