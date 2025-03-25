package gin

import (
	"net/http"
	"testing"

	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// 测试 noneWrappedData 函数 - 处理不同类型的数据
func Test_noneWrappedData2(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// 创建模拟的 XContext
	mockContext := NewMockXContext(controller)

	// 测试场景1: nil 数据
	mockContext.EXPECT().String(http.StatusOK, "").Times(1)
	noneWrappedData(mockContext, nil, http.StatusOK)

	// 测试场景2: 错误类型
	testErr := gone.NewBusinessError("test error", 400, nil)
	mockContext.EXPECT().String(http.StatusOK, testErr.Error()).Times(1)
	noneWrappedData(mockContext, testErr, http.StatusOK)

	// 测试场景3: 结构体类型
	testStruct := struct {
		Name string
		Age  int
	}{"John", 30}
	mockContext.EXPECT().JSON(http.StatusOK, testStruct).Times(1)
	noneWrappedData(mockContext, testStruct, http.StatusOK)

	// 测试场景4: 指针类型 - 结构体指针
	testStructPtr := &testStruct
	mockContext.EXPECT().JSON(http.StatusOK, testStructPtr).Times(1)
	noneWrappedData(mockContext, testStructPtr, http.StatusOK)

	// 测试场景5: 指针类型 - 基本类型指针
	testInt := 123
	testIntPtr := &testInt
	mockContext.EXPECT().String(http.StatusOK, "123").Times(1)
	noneWrappedData(mockContext, testIntPtr, http.StatusOK)

	// 测试场景6: 基本类型 - 整数
	mockContext.EXPECT().String(http.StatusOK, "123").Times(1)
	noneWrappedData(mockContext, 123, http.StatusOK)

	// 测试场景7: 基本类型 - 字符串
	mockContext.EXPECT().String(http.StatusOK, "test string").Times(1)
	noneWrappedData(mockContext, "test string", http.StatusOK)

	// 测试场景8: 切片类型
	testSlice := []string{"a", "b", "c"}
	mockContext.EXPECT().JSON(http.StatusOK, testSlice).Times(1)
	noneWrappedData(mockContext, testSlice, http.StatusOK)

	// 测试场景9: 映射类型
	testMap := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
	}
	mockContext.EXPECT().JSON(http.StatusOK, testMap).Times(1)
	noneWrappedData(mockContext, testMap, http.StatusOK)
}

// 测试 responser.innerErrorMsg 方法
func Test_responser_innerErrorMsg(t *testing.T) {
	// 测试场景1: 显示内部错误详情
	r1 := &responser{
		doNotShowInnerErrorDetail: false,
	}
	iErr := gone.NewInnerError("test inner error", 500)
	msg := r1.innerErrorMsg(iErr.(gone.InnerError))
	assert.Contains(t, msg, "test inner error")

	// 测试场景2: 不显示内部错误详情
	r2 := &responser{
		doNotShowInnerErrorDetail: true,
	}
	msg = r2.innerErrorMsg(iErr.(gone.InnerError))
	assert.Equal(t, InternalServerError, msg)
}

// 测试 responser.processChan 方法 - 处理不同类型的通道数据
func Test_responser_processChan_DifferentTypes(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// 创建模拟的 Logger
	mockLogger := NewMockLogger(controller)

	// 创建 responser 实例
	r := &responser{
		Logger:                    mockLogger,
		doNotShowInnerErrorDetail: true,
	}

	// 创建测试 ResponseWriter
	w := NewMockResponseWriter(controller)

	var response string
	w.EXPECT().Header().Return(http.Header{}).AnyTimes()
	w.EXPECT().Flush().AnyTimes()
	w.EXPECT().
		WriteString(gomock.Any()).
		DoAndReturn(func(data string) (int, error) {
			response += data
			return len(data), nil
		}).
		AnyTimes()
	w.EXPECT().CloseNotify().AnyTimes()

	// 测试场景1: 处理 InnerError 类型
	ch1 := make(chan gone.InnerError, 1)
	ch1 <- gone.NewInnerError("test inner error", 500).(gone.InnerError)
	close(ch1)

	// 设置期望
	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()

	// 处理通道
	r.processChan(ch1, w)

	// 验证响应包含预期的数据
	assert.Contains(t, response, "event: data")
	assert.Contains(t, response, `"code":500`)
	assert.Contains(t, response, `"msg":"Internal Server Error"`)
	assert.Contains(t, response, "event: done")

	// 重置响应
	response = ""

	// 测试场景2: 处理 BusinessError 类型
	ch2 := make(chan gone.BusinessError, 1)
	ch2 <- gone.NewBusinessError("test business error", 400, "error data")
	close(ch2)

	// 处理通道
	r.processChan(ch2, w)

	// 验证响应包含预期的数据
	assert.Contains(t, response, "event: data")
	assert.Contains(t, response, `"code":400`)
	assert.Contains(t, response, `"msg":"test business error"`)
	assert.Contains(t, response, `"data":"error data"`)
	assert.Contains(t, response, "event: done")

	// 重置响应
	response = ""

	// 测试场景3: 处理 Error 类型
	ch3 := make(chan gone.Error, 1)
	ch3 <- gone.NewError(400, "test error", 400)
	close(ch3)

	// 处理通道
	r.processChan(ch3, w)

	// 验证响应包含预期的数据
	assert.Contains(t, response, "event: data")
	assert.Contains(t, response, `"code":400`)
	assert.Contains(t, response, `"msg":"test error"`)
	assert.Contains(t, response, "event: done")

	// 重置响应
	response = ""

	// 测试场景4: 处理普通错误类型
	ch4 := make(chan error, 1)
	ch4 <- gone.NewError(400, "test error", 400)
	close(ch4)

	// 处理通道
	r.processChan(ch4, w)

	// 验证响应包含预期的数据
	assert.Contains(t, response, "event: data")
	assert.Contains(t, response, `"code":400`)
	assert.Contains(t, response, `"msg":"test error"`)
	assert.Contains(t, response, "event: done")

	// 重置响应
	response = ""

	// 测试场景5: 处理普通数据类型
	ch5 := make(chan string, 1)
	ch5 <- "test data"
	close(ch5)

	// 处理通道
	r.processChan(ch5, w)

	// 验证响应包含预期的数据
	assert.Contains(t, response, "event: data")
	assert.Contains(t, response, `"test data"`)
	assert.Contains(t, response, "event: done")
}

// 测试 responser.processChan 方法 - 处理 SSE 写入错误
func Test_responser_processChan_WriteError(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// 创建模拟的 Logger
	mockLogger := NewMockLogger(controller)

	// 创建 responser 实例
	r := &responser{
		Logger: mockLogger,
	}

	// 创建测试 ResponseWriter，模拟写入错误
	w := NewMockResponseWriter(controller)

	w.EXPECT().Header().Return(http.Header{}).AnyTimes()
	w.EXPECT().Flush().AnyTimes()
	w.EXPECT().
		WriteString(gomock.Any()).
		Return(0, http.ErrBodyNotAllowed).
		AnyTimes()
	w.EXPECT().CloseNotify().AnyTimes()

	// 创建测试通道
	ch := make(chan string, 1)
	ch <- "test data"
	close(ch)

	// 设置期望 - 应该记录错误
	mockLogger.EXPECT().Errorf("write data error: %v", http.ErrBodyNotAllowed).Times(1)
	mockLogger.EXPECT().Errorf("write 'end' error: %v", http.ErrBodyNotAllowed).Times(1)

	// 处理通道
	r.processChan(ch, w)
}
