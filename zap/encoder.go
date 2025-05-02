package gone_zap

import (
	"github.com/gone-io/goner/g"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type traceEncoder struct {
	zapcore.Encoder
	tracer g.Tracer
}

const traceIdKey = "traceId"

func (e *traceEncoder) EncodeEntry(entry zapcore.Entry, fields []Field) (*buffer.Buffer, error) {
	traceId := e.tracer.GetTraceId()
	if traceId != "" {
		fields = append(fields, zap.String(traceIdKey, traceId))
	}
	return e.Encoder.EncodeEntry(entry, fields)
}

func NewTraceEncoder(encoder zapcore.Encoder, tracer g.Tracer) zapcore.Encoder {
	return &traceEncoder{
		Encoder: encoder,
		tracer:  tracer,
	}
}
