package gone_zap

//
//// 测试wrappedLogger的各种方法
//func TestWrappedLogger(t *testing.T) {
//	// 创建一个基础的zap.Logger用于测试
//	config := zap.NewDevelopmentConfig()
//	logger, err := config.Build()
//	assert.Nil(t, err)
//
//	// 创建wrappedLogger实例
//	wrapped := &wrappedLogger{Logger: logger}
//
//	// 测试sugar方法
//	sugar := wrapped.sugar()
//	assert.NotNil(t, sugar)
//
//	// 测试Named方法 - 空名称
//	named := wrapped.Named("")
//	assert.Equal(t, wrapped, named)
//
//	// 测试Named方法 - 有名称
//	named = wrapped.Named("test-logger")
//	assert.NotEqual(t, wrapped, named)
//	assert.IsType(t, &wrappedLogger{}, named)
//
//	// 测试WithOptions方法 - 空选项
//	withOpts := wrapped.WithOptions()
//	assert.Equal(t, wrapped, withOpts)
//
//	// 测试WithOptions方法 - 有选项
//	withOpts = wrapped.WithOptions(zap.AddCaller())
//	assert.IsType(t, &wrappedLogger{}, withOpts)
//
//	// 测试With方法 - 空字段
//	withFields := wrapped.With()
//	assert.Equal(t, wrapped, withFields)
//
//	// 测试With方法 - 有字段
//	withFields = wrapped.With(zap.String("key", "value"))
//	assert.NotEqual(t, wrapped, withFields)
//	assert.IsType(t, &wrappedLogger{}, withFields)
//
//	// 测试Sugar方法
//	sugarLogger := wrapped.Sugar()
//	assert.NotNil(t, sugarLogger)
//}
//
//// 测试wrappedLogger的日志方法
//func TestWrappedLogger_LogMethods(t *testing.T) {
//	// 使用内存缓冲区捕获日志输出
//	core, recorded := observer.New(zapcore.InfoLevel)
//	logger := zap.New(core)
//
//	// 创建wrappedLogger实例
//	wrapped := &wrappedLogger{Logger: logger}
//
//	// 测试各种日志级别方法
//	wrapped.Debug("debug message", zap.String("key", "value"))
//	wrapped.Info("info message", zap.String("key", "value"))
//	wrapped.Warn("warn message", zap.String("key", "value"))
//	wrapped.Error("error message", zap.String("key", "value"))
//
//	// 验证日志输出
//	logs := recorded.All()
//	assert.Equal(t, 3, len(logs)) // Debug级别低于Info，不会被记录
//
//	// 验证Info日志
//	assert.Equal(t, "info message", logs[0].Message)
//	assert.Equal(t, zapcore.InfoLevel, logs[0].Level)
//	assert.Equal(t, "value", logs[0].ContextMap()["key"])
//
//	// 验证Warn日志
//	assert.Equal(t, "warn message", logs[1].Message)
//	assert.Equal(t, zapcore.WarnLevel, logs[1].Level)
//
//	// 验证Error日志
//	assert.Equal(t, "error message", logs[2].Message)
//	assert.Equal(t, zapcore.ErrorLevel, logs[2].Level)
//}
