package main

import (
	"github.com/gone-io/gone/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

var _ zapcore.Encoder = (*UseCustomerEncoder)(nil)

// 演示如何使用自定义的Encoder并加载到gone中
// func init() {
// 	gone.Load(NewUseCustomerEncoder())
// }

func NewUseCustomerEncoder() *UseCustomerEncoder {
	return &UseCustomerEncoder{
		Encoder: zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig()),
	}
}

type UseCustomerEncoder struct {
	zapcore.Encoder
	gone.Flag
}

func (e *UseCustomerEncoder) EncodeEntry(entry zapcore.Entry, fields []zap.Field) (*buffer.Buffer, error) {
	//do something
	return e.Encoder.EncodeEntry(entry, fields)
}
