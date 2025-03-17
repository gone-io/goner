package tracer

import (
	"sync"
	"testing"
)

// 基准测试 tracer.SetTraceId 方法
func BenchmarkTracer_SetTraceId(b *testing.B) {
	tracer := &tracer{}
	customTraceId := "custom-trace-id"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tracer.SetTraceId(customTraceId, func() {
			// 空函数，只测试SetTraceId的性能
		})
	}
}

// 基准测试 tracerOverGid.SetTraceId 方法
func BenchmarkTracerOverGid_SetTraceId(b *testing.B) {
	tracer := &tracerOverGid{}
	customTraceId := "custom-trace-id"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tracer.SetTraceId(customTraceId, func() {
			// 空函数，只测试SetTraceId的性能
		})
	}
}

// 基准测试 tracer.GetTraceId 方法
func BenchmarkTracer_GetTraceId(b *testing.B) {
	tracer := &tracer{}
	customTraceId := "custom-trace-id"

	// 先设置一个traceId
	tracer.SetTraceId(customTraceId, func() {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = tracer.GetTraceId()
		}
	})
}

// 基准测试 tracerOverGid.GetTraceId 方法
func BenchmarkTracerOverGid_GetTraceId(b *testing.B) {
	tracer := &tracerOverGid{}
	customTraceId := "custom-trace-id"

	// 先设置一个traceId
	tracer.SetTraceId(customTraceId, func() {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = tracer.GetTraceId()
		}
	})
}

// 基准测试 tracer.Go 方法
func BenchmarkTracer_Go(b *testing.B) {
	tracer := &tracer{}
	customTraceId := "custom-trace-id"
	var wg sync.WaitGroup

	tracer.SetTraceId(customTraceId, func() {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			wg.Add(1)
			tracer.Go(func() {
				_ = tracer.GetTraceId()
				wg.Done()
			})
		}
	})

	wg.Wait()
}

// 基准测试 tracerOverGid.Go 方法
func BenchmarkTracerOverGid_Go(b *testing.B) {
	tracer := &tracerOverGid{}
	customTraceId := "custom-trace-id"
	var wg sync.WaitGroup

	tracer.SetTraceId(customTraceId, func() {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			wg.Add(1)
			tracer.Go(func() {
				_ = tracer.GetTraceId()
				wg.Done()
			})
		}
	})

	wg.Wait()
}

// 基准测试 tracer 在并发环境下的性能
func BenchmarkTracer_Concurrent(b *testing.B) {
	tracer := &tracer{}
	customTraceId := "custom-trace-id"
	const goroutines = 10

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(goroutines)

		for j := 0; j < goroutines; j++ {
			go func() {
				defer wg.Done()
				tracer.SetTraceId(customTraceId, func() {
					_ = tracer.GetTraceId()
				})
			}()
		}

		wg.Wait()
	}
}

// 基准测试 tracerOverGid 在并发环境下的性能
func BenchmarkTracerOverGid_Concurrent(b *testing.B) {
	tracer := &tracerOverGid{}
	customTraceId := "custom-trace-id"
	const goroutines = 10

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(goroutines)

		for j := 0; j < goroutines; j++ {
			go func() {
				defer wg.Done()
				tracer.SetTraceId(customTraceId, func() {
					_ = tracer.GetTraceId()
				})
			}()
		}

		wg.Wait()
	}
}

// 基准测试 tracer 在嵌套调用场景下的性能
func BenchmarkTracer_Nested(b *testing.B) {
	tracer := &tracer{}
	customTraceId := "custom-trace-id"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tracer.SetTraceId(customTraceId, func() {
			_ = tracer.GetTraceId()
			tracer.SetTraceId("inner-"+customTraceId, func() {
				_ = tracer.GetTraceId()
			})
		})
	}
}

// 基准测试 tracerOverGid 在嵌套调用场景下的性能
func BenchmarkTracerOverGid_Nested(b *testing.B) {
	tracer := &tracerOverGid{}
	customTraceId := "custom-trace-id"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tracer.SetTraceId(customTraceId, func() {
			_ = tracer.GetTraceId()
			tracer.SetTraceId("inner-"+customTraceId, func() {
				_ = tracer.GetTraceId()
			})
		})
	}
}
