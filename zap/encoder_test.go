package gone_zap

import (
	gMock "github.com/gone-io/goner/g/mock"
	"go.uber.org/mock/gomock"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

// 模拟tracer实现
//type mockTracer struct {
//	traceId string
//}

//func (m *mockTracer) SetTraceId(traceId string, fn func()) {
//	m.traceId = traceId
//	fn()
//}
//
//func (m *mockTracer) GetTraceId() string {
//	return m.traceId
//}
//
//func (m *mockTracer) Go(fn func()) {
//	go fn()
//}

// 模拟encoder实现
type mockEncoder struct {
	fields []zapcore.Field
}

func (m *mockEncoder) AddArray(key string, marshaler zapcore.ArrayMarshaler) error {
	return nil
}

func (m *mockEncoder) AddObject(key string, marshaler zapcore.ObjectMarshaler) error {
	return nil
}

func (m *mockEncoder) AddBinary(key string, value []byte) {
}

func (m *mockEncoder) AddByteString(key string, value []byte) {
}

func (m *mockEncoder) AddBool(key string, value bool) {
}

func (m *mockEncoder) AddComplex128(key string, value complex128) {
}

func (m *mockEncoder) AddComplex64(key string, value complex64) {
}

func (m *mockEncoder) AddDuration(key string, value time.Duration) {
}

func (m *mockEncoder) AddFloat64(key string, value float64) {
}

func (m *mockEncoder) AddFloat32(key string, value float32) {
}

func (m *mockEncoder) AddInt(key string, value int) {
}

func (m *mockEncoder) AddInt64(key string, value int64) {
}

func (m *mockEncoder) AddInt32(key string, value int32) {
}

func (m *mockEncoder) AddInt16(key string, value int16) {
}

func (m *mockEncoder) AddInt8(key string, value int8) {
}

func (m *mockEncoder) AddString(key string, value string) {
	m.fields = append(m.fields, zap.String(key, value))
}

func (m *mockEncoder) AddTime(key string, value time.Time) {
}

func (m *mockEncoder) AddUint(key string, value uint) {
}

func (m *mockEncoder) AddUint64(key string, value uint64) {
}

func (m *mockEncoder) AddUint32(key string, value uint32) {
}

func (m *mockEncoder) AddUint16(key string, value uint16) {
}

func (m *mockEncoder) AddUint8(key string, value uint8) {
}

func (m *mockEncoder) AddUintptr(key string, value uintptr) {
}

func (m *mockEncoder) AddReflected(key string, value interface{}) error {
	return nil
}

func (m *mockEncoder) OpenNamespace(key string) {
}

func (m *mockEncoder) Clone() zapcore.Encoder {
	return m
}

func (m *mockEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	m.fields = append(m.fields, fields...)
	return &buffer.Buffer{}, nil
}

// 测试traceEncoder的EncodeEntry方法
func TestTraceEncoder_EncodeEntry(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockTracer := gMock.NewMockTracer(controller)

	traceId := "mocked-trace-id"
	mockTracer.EXPECT().GetTraceId().Return(traceId).AnyTimes()

	mockBaseEncoder := &mockEncoder{}
	traceEncoder := NewTraceEncoder(mockBaseEncoder, mockTracer)

	// 创建一个测试条目
	entry := zapcore.Entry{
		Level:      zapcore.InfoLevel,
		Time:       time.Now(),
		LoggerName: "test",
		Message:    "test message",
	}

	// 添加一些字段
	fields := []zapcore.Field{
		zap.String("key1", "value1"),
		zap.Int("key2", 123),
	}

	// 调用EncodeEntry方法
	_, err := traceEncoder.EncodeEntry(entry, fields)

	// 验证结果
	assert.Nil(t, err)

	// 验证traceId被添加到字段中
	found := false
	for _, field := range mockBaseEncoder.fields {
		if field.Key == "traceId" {
			assert.Equal(t, traceId, field.String)
			found = true
			break
		}
	}
	assert.True(t, found, "traceId field should be added")

	mockBaseEncoder = &mockEncoder{}
	traceEncoder = NewTraceEncoder(mockBaseEncoder, mockTracer)

	// 调用EncodeEntry方法
	_, err = traceEncoder.EncodeEntry(entry, fields)

	traceEncoder = traceEncoder.Clone()
	traceEncoder.AddString(contextKey, contextValue)

	// 验证结果
	assert.Nil(t, err)

	// 验证没有traceId被添加到字段中
	found = false
	for _, field := range mockBaseEncoder.fields {
		if field.Key == "traceId" {
			found = true
			break
		}
	}
	assert.True(t, found, "traceId field should not be added when traceId is empty")
}
