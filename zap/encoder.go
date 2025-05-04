package gone_zap

import (
	"github.com/gone-io/goner/g"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type traceEncoder struct {
	zapcore.Encoder
	tracer       g.Tracer
	addedTraceId bool
}

const traceIdKey = "traceId"
const contextKey = "context"
const contextValue = "context.Background.WithValue(trace.traceContextKeyType, *trace.recordingSpan)"

func (e *traceEncoder) EncodeEntry(entry zapcore.Entry, fields []Field) (*buffer.Buffer, error) {
	if e.tracer != nil && !e.addedTraceId {
		traceId := e.tracer.GetTraceId()
		if traceId != "" {
			fields = append(fields, zap.String(traceIdKey, traceId))
		}
	}

	return e.Encoder.EncodeEntry(entry, fields)
}

func (e *traceEncoder) AddString(key, value string) {
	if key == traceIdKey && value != "" {
		e.addedTraceId = true
	}
	if key == contextKey && value == contextValue {
		return
	}
	e.Encoder.AddString(key, value)
}

func (e *traceEncoder) Clone() zapcore.Encoder {
	return &traceEncoder{
		Encoder: e.Encoder.Clone(),
		tracer:  e.tracer,
	}
}

func NewTraceEncoder(encoder zapcore.Encoder, tracer g.Tracer) zapcore.Encoder {
	return &traceEncoder{
		Encoder: encoder,
		tracer:  tracer,
	}
}
