package gin

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// 测试 NewGinResponser 函数
func TestNewGinResponser(t *testing.T) {
	r := NewGinResponser()
	assert.NotNil(t, r)
	assert.IsType(t, &responser{}, r)
}

// 测试 SetWrappedDataFunc 方法
func Test_responser_SetWrappedDataFunc(t *testing.T) {
	r := &responser{
		wrappedDataFunc: wrapFunc,
	}

	// 自定义包装函数
	customWrapFunc := func(code int, msg string, data any) any {
		return map[string]any{
			"status":  code,
			"message": msg,
			"result":  data,
		}
	}

	// 设置自定义包装函数
	r.SetWrappedDataFunc(customWrapFunc)

	// 验证包装函数已被设置
	result := r.wrappedDataFunc(200, "success", "test")
	expected := map[string]any{
		"status":  200,
		"message": "success",
		"result":  "test",
	}

	assert.Equal(t, expected, result)
}

// 测试 Success 方法 - 包装数据模式
func Test_responser_Success_WrappedData(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// 创建模拟的 XContext
	mockContext := NewMockXContext(controller)

	// 创建 responser 实例
	r := &responser{
		wrappedDataFunc:   wrapFunc,
		returnWrappedData: true,
	}

	// 测试普通数据
	mockContext.EXPECT().JSON(http.StatusOK, gomock.Any()).Do(func(status int, data any) {
		resData, ok := data.(*res[any])
		assert.True(t, ok)
		assert.Equal(t, 0, resData.Code)
		assert.Equal(t, "", resData.Msg)
		assert.Equal(t, "test data", resData.Data)
	})

	r.Success(mockContext, "test data")

	// 测试 BusinessError 类型
	bErr := gone.NewBusinessError("business error", 400, "error data")
	mockContext.EXPECT().JSON(http.StatusOK, gomock.Any()).Do(func(status int, data any) {
		resData, ok := data.(*res[any])
		assert.True(t, ok)
		assert.Equal(t, 400, resData.Code)
		assert.Equal(t, "business error", resData.Msg)
		assert.Equal(t, "error data", resData.Data)
	})

	r.Success(mockContext, bErr)
}

// 测试 Success 方法 - 非包装数据模式
func Test_responser_Success_NonWrappedData(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// 创建模拟的 XContext
	mockContext := NewMockXContext(controller)

	// 创建 responser 实例
	r := &responser{
		wrappedDataFunc:   wrapFunc,
		returnWrappedData: false,
	}

	// 测试 nil 数据
	mockContext.EXPECT().String(http.StatusOK, "").Times(1)
	r.Success(mockContext, nil)

	// 测试错误类型
	testErr := errors.New("test error")
	mockContext.EXPECT().String(http.StatusOK, "test error").Times(1)
	r.Success(mockContext, testErr)

	// 测试结构体类型
	testStruct := struct {
		Name string
		Age  int
	}{"John", 30}
	mockContext.EXPECT().JSON(http.StatusOK, testStruct).Times(1)
	r.Success(mockContext, testStruct)

	// 测试指针类型
	testPtr := &testStruct
	mockContext.EXPECT().JSON(http.StatusOK, testPtr).Times(1)
	r.Success(mockContext, testPtr)

	// 测试基本类型
	mockContext.EXPECT().String(http.StatusOK, "123").Times(1)
	r.Success(mockContext, 123)
}

// 测试 Failed 方法 - 包装数据模式
func Test_responser_Failed_WrappedData(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// 创建模拟的 XContext 和 Logger
	mockContext := NewMockXContext(controller)
	mockLogger := NewMockLogger(controller)

	// 创建 responser 实例
	r := &responser{
		Logger:                    mockLogger,
		wrappedDataFunc:           wrapFunc,
		returnWrappedData:         true,
		doNotShowInnerErrorDetail: true,
	}

	// 测试 nil 错误
	mockContext.EXPECT().JSON(http.StatusBadRequest, gomock.Any()).Do(func(status int, data any) {
		resData, ok := data.(*res[any])
		assert.True(t, ok)
		assert.Equal(t, 0, resData.Code)
		assert.Equal(t, "", resData.Msg)
		assert.Nil(t, resData.Data)
	})

	r.Failed(mockContext, nil)

	// 测试 BusinessError 类型
	bErr := gone.NewBusinessError("business error", 400, "error data")
	mockContext.EXPECT().JSON(http.StatusOK, gomock.Any()).Do(func(status int, data any) {
		resData, ok := data.(*res[any])
		assert.True(t, ok)
		assert.Equal(t, 400, resData.Code)
		assert.Equal(t, "business error", resData.Msg)
		assert.Equal(t, "error data", resData.Data)
	})

	r.Failed(mockContext, bErr)

	// 测试 InnerError 类型
	iErr := gone.NewInnerError("inner error", 500)
	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1)
	mockContext.EXPECT().JSON(http.StatusInternalServerError, gomock.Any()).Do(func(status int, data any) {
		resData, ok := data.(*res[any])
		assert.True(t, ok)
		assert.Equal(t, 500, resData.Code)
		assert.Equal(t, "Internal Server Error", resData.Msg)
		assert.Nil(t, resData.Data)
	})

	r.Failed(mockContext, iErr)

	// 测试普通错误
	pErr := gone.NewParameterError("parameter error")
	mockContext.EXPECT().JSON(http.StatusBadRequest, gomock.Any()).Do(func(status int, data any) {
		resData, ok := data.(*res[any])
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, resData.Code)
		assert.Equal(t, "parameter error", resData.Msg)
		assert.Nil(t, resData.Data)
	})

	r.Failed(mockContext, pErr)
}

// 测试 Failed 方法 - 非包装数据模式
func Test_responser_Failed_NonWrappedData(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// 创建模拟的 XContext 和 Logger
	mockContext := NewMockXContext(controller)
	mockLogger := NewMockLogger(controller)

	// 创建 responser 实例
	r := &responser{
		Logger:                    mockLogger,
		wrappedDataFunc:           wrapFunc,
		returnWrappedData:         false,
		doNotShowInnerErrorDetail: true,
	}

	// 测试 nil 错误
	mockContext.EXPECT().String(http.StatusBadRequest, "").Times(1)
	r.Failed(mockContext, nil)

	// 测试 InnerError 类型
	iErr := gone.NewInnerError("inner error", 500)
	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1)
	mockContext.EXPECT().String(http.StatusInternalServerError, InternalServerError).Times(1)
	r.Failed(mockContext, iErr)

	// 测试普通错误
	pErr := gone.NewParameterError("parameter error")
	mockContext.EXPECT().String(http.StatusBadRequest, "GoneError(code=400); parameter error").Times(1)
	r.Failed(mockContext, pErr)
}

// 测试 ProcessResults 方法 - 处理错误
func Test_responser_ProcessResults_Error(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// 创建模拟的 XContext, ResponseWriter 和 Logger
	mockContext := NewMockXContext(controller)
	mockWriter := NewMockResponseWriter(controller)
	mockLogger := NewMockLogger(controller)

	// 创建 responser 实例
	r := &responser{
		Logger:            mockLogger,
		wrappedDataFunc:   wrapFunc,
		returnWrappedData: true,
	}

	mockWriter.EXPECT().Written().Times(1)

	// 测试处理错误
	testErr := errors.New("test error")
	mockContext.EXPECT().JSON(gomock.Any(), gomock.Any()).Times(1)
	mockContext.EXPECT().Abort().Times(1)
	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

	r.ProcessResults(mockContext, mockWriter, true, "testFunc", testErr)
}

// 测试 ProcessResults 方法 - 处理通道
func Test_responser_ProcessResults_Channel(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// 创建模拟的 XContext 和 Logger
	mockContext := NewMockXContext(controller)
	mockLogger := NewMockLogger(controller)

	// 创建 responser 实例
	r := &responser{
		Logger:            mockLogger,
		wrappedDataFunc:   wrapFunc,
		returnWrappedData: true,
	}

	// 创建测试通道
	ch := make(chan string, 2)
	ch <- "test1"
	ch <- "test2"
	close(ch)

	// 创建测试 ResponseWriter
	w := NewMockResponseWriter(controller)

	var response string

	w.EXPECT().Written().AnyTimes()
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

	// 处理通道
	r.ProcessResults(mockContext, w, true, "testFunc", ch)

	// 验证响应
	assert.Contains(t, response, "event: data")
	assert.Contains(t, response, `data: "test1"`)
	assert.Contains(t, response, `data: "test2"`)
	assert.Contains(t, response, "event: done")
}

// 测试 ProcessResults 方法 - 处理 io.Reader
func Test_responser_ProcessResults_Reader(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// 创建模拟的 XContext 和 Logger
	mockContext := NewMockXContext(controller)
	mockLogger := NewMockLogger(controller)

	// 创建 responser 实例
	r := &responser{
		Logger:            mockLogger,
		wrappedDataFunc:   wrapFunc,
		returnWrappedData: true,
	}

	// 创建测试 Reader
	reader := strings.NewReader("test reader data")

	// 创建测试 ResponseWriter
	w := NewMockResponseWriter(controller)

	var response string
	w.EXPECT().Written().AnyTimes()
	w.EXPECT().
		WriteString(gomock.Any()).
		DoAndReturn(func(data string) (int, error) {
			response += data
			return len(data), nil
		}).
		AnyTimes()

	// 处理 Reader
	r.ProcessResults(mockContext, w, true, "testFunc", reader)

	// 验证响应
	assert.Equal(t, "test reader data", response)
}

// 测试 ProcessResults 方法 - 处理普通数据
func Test_responser_ProcessResults_NormalData(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// 创建模拟的 XContext 和 Logger
	mockContext := NewMockXContext(controller)
	mockLogger := NewMockLogger(controller)

	// 创建 responser 实例
	r := &responser{
		Logger:            mockLogger,
		wrappedDataFunc:   wrapFunc,
		returnWrappedData: true,
	}

	// 创建测试 ResponseWriter
	w := NewMockResponseWriter(controller)

	isWritten := false
	w.EXPECT().
		Written().
		DoAndReturn(func() bool {
			return isWritten
		}).
		Times(1)
	// 测试处理普通数据
	mockContext.EXPECT().
		JSON(http.StatusOK, gomock.Any()).
		Do(func(httpStatus int, data any) {
			isWritten = true
		}).
		Times(1)

	r.ProcessResults(mockContext, w, true, "testFunc", "test data")

	isWritten = false
	// 测试处理 nil 数据
	w.EXPECT().
		Written().
		DoAndReturn(func() bool {
			return isWritten
		}).
		Times(1)
	// 测试处理普通数据
	mockContext.EXPECT().
		JSON(http.StatusOK, gomock.Any()).
		Do(func(httpStatus int, data any) {
			isWritten = true
		}).
		Times(1)
	r.ProcessResults(mockContext, w, true, "testFunc", nil)

	isWritten = false
	// 测试处理多个结果
	w.EXPECT().
		Written().
		DoAndReturn(func() bool {
			return isWritten
		}).
		Times(1)
	// 测试处理普通数据
	mockContext.EXPECT().
		JSON(http.StatusOK, gomock.Any()).
		Do(func(httpStatus int, data any) {
			isWritten = true
		}).
		Times(1)
	r.ProcessResults(mockContext, w, true, "testFunc", nil, "test data")
}

// 测试 processChan 方法
func Test_responser_processChan(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// 创建模拟的 Logger
	mockLogger := NewMockLogger(controller)

	// 创建 responser 实例
	r := &responser{
		Logger: mockLogger,
	}

	// 创建测试通道 - 正常数据
	ch1 := make(chan string, 2)
	ch1 <- "test1"
	ch1 <- "test2"
	close(ch1)

	// 创建测试 ResponseWriter
	w1 := NewMockResponseWriter(controller)
	w1.EXPECT().Header().Return(http.Header{}).AnyTimes()
	w1.EXPECT().Flush().AnyTimes()

	var response1 string
	w1.EXPECT().
		WriteString(gomock.Any()).
		DoAndReturn(func(data string) (int, error) {
			response1 += data
			return len(data), nil
		}).
		AnyTimes()
	w1.EXPECT().CloseNotify().Times(1)

	// 处理通道
	r.processChan(ch1, w1)

	// 验证响应
	assert.Contains(t, response1, "event: data")
	assert.Contains(t, response1, `data: "test1"`)
	assert.Contains(t, response1, `data: "test2"`)
	assert.Contains(t, response1, "event: done")

	// 创建测试通道 - 错误数据
	ch2 := make(chan any, 2)
	ch2 <- errors.New("test error")
	ch2 <- "normal data"
	close(ch2)

	// 创建测试 ResponseWriter
	w2 := NewMockResponseWriter(controller)
	w2.EXPECT().Header().Return(http.Header{}).AnyTimes()
	w2.EXPECT().Flush().AnyTimes()

	var response2 string
	w2.EXPECT().
		WriteString(gomock.Any()).
		DoAndReturn(func(data string) (int, error) {
			response2 += data
			return len(data), nil
		}).
		AnyTimes()
	w2.EXPECT().CloseNotify().Times(1)

	// 设置 Logger 的期望行为
	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()

	// 处理通道
	r.processChan(ch2, w2)

	// 验证响应
	assert.Contains(t, response2, "event: data")
	assert.Contains(t, response2, `"code":500`)
	assert.Contains(t, response2, `test error`)
	assert.Contains(t, response2, `data: "normal data"`)
	assert.Contains(t, response2, "event: done")
}

// 测试 wrapFunc 函数
func Test_wrapFunc(t *testing.T) {
	// 测试正常数据
	result1 := wrapFunc(0, "", "test data")
	resData1, ok := result1.(*res[any])
	assert.True(t, ok)
	assert.Equal(t, 0, resData1.Code)
	assert.Equal(t, "", resData1.Msg)
	assert.Equal(t, "test data", resData1.Data)

	// 测试错误数据
	result2 := wrapFunc(400, "error message", nil)
	resData2, ok := result2.(*res[any])
	assert.True(t, ok)
	assert.Equal(t, 400, resData2.Code)
	assert.Equal(t, "error message", resData2.Msg)
	assert.Nil(t, resData2.Data)
}

// 测试 noneWrappedData 函数
func Test_noneWrappedData(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	// 创建模拟的 XContext
	mockContext := NewMockXContext(controller)

	// 测试 nil 数据
	mockContext.EXPECT().String(http.StatusOK, "").Times(1)
	noneWrappedData(mockContext, nil, http.StatusOK)

	// 测试错误类型
	testErr := errors.New("test error")
	mockContext.EXPECT().String(http.StatusOK, "test error").Times(1)
	noneWrappedData(mockContext, testErr, http.StatusOK)

	// 测试结构体类型
	testStruct := struct {
		Name string
		Age  int
	}{"John", 30}
	mockContext.EXPECT().JSON(http.StatusOK, testStruct).Times(1)
	noneWrappedData(mockContext, testStruct, http.StatusOK)

	// 测试 Map 类型
	testMap := map[string]any{"name": "John", "age": 30}
	mockContext.EXPECT().JSON(http.StatusOK, testMap).Times(1)
	noneWrappedData(mockContext, testMap, http.StatusOK)

	// 测试 Slice 类型
	testSlice := []string{"a", "b", "c"}
	mockContext.EXPECT().JSON(http.StatusOK, testSlice).Times(1)
	noneWrappedData(mockContext, testSlice, http.StatusOK)

	// 测试指针类型
	testPtr := &testStruct
	mockContext.EXPECT().JSON(http.StatusOK, testPtr).Times(1)
	noneWrappedData(mockContext, testPtr, http.StatusOK)

	// 测试基本类型
	mockContext.EXPECT().String(http.StatusOK, "123").Times(1)
	noneWrappedData(mockContext, 123, http.StatusOK)

	// 测试指向基本类型的指针
	testInt := 123
	testIntPtr := &testInt
	mockContext.EXPECT().String(http.StatusOK, "123").Times(1)
	noneWrappedData(mockContext, testIntPtr, http.StatusOK)
}
