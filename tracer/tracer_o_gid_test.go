package tracer

import (
	"sync"
	"testing"

	"github.com/petermattis/goid"
)

func Test_GetGoId(t *testing.T) {
	gid := goid.Get()
	if gid == 0 {
		t.Fatal("can not get goid")
	}
	t.Logf("gid=%d", gid)
}

func TestTracerOverGid_SetAndGetTraceId(t *testing.T) {
	tracer := &tracerOverGid{}

	// 测试设置自定义traceId
	customTraceId := "custom-trace-id"
	var gotTraceId string

	tracer.SetTraceId(customTraceId, func() {
		gotTraceId = tracer.GetTraceId()
	})

	if gotTraceId != customTraceId {
		t.Errorf("SetTraceId(%s) = %s; want %s", customTraceId, gotTraceId, customTraceId)
	}

	// 测试空traceId时自动生成
	var generatedTraceId string
	tracer.SetTraceId("", func() {
		generatedTraceId = tracer.GetTraceId()
	})

	if generatedTraceId == "" {
		t.Error("SetTraceId(\"\") should generate a non-empty traceId")
	}

	// 测试在回调函数外部获取traceId
	outsideTraceId := tracer.GetTraceId()
	if outsideTraceId != "" {
		t.Errorf("GetTraceId() outside of SetTraceId callback = %s; want empty string", outsideTraceId)
	}
}

func TestTracerOverGid_Go(t *testing.T) {
	tracer := &tracerOverGid{}
	customTraceId := "go-routine-trace-id"

	var wg sync.WaitGroup
	var childTraceId string

	tracer.SetTraceId(customTraceId, func() {
		parentTraceId := tracer.GetTraceId()
		if parentTraceId != customTraceId {
			t.Errorf("Parent goroutine traceId = %s; want %s", parentTraceId, customTraceId)
		}

		wg.Add(1)
		tracer.Go(func() {
			defer wg.Done()
			childTraceId = tracer.GetTraceId()
		})
	})

	wg.Wait()

	if childTraceId != customTraceId {
		t.Errorf("Child goroutine traceId = %s; want %s", childTraceId, customTraceId)
	}
}

func TestTracerOverGid_MultipleGoroutines(t *testing.T) {
	tracer := &tracerOverGid{}
	customTraceId := "multi-goroutine-trace-id"

	var wg sync.WaitGroup
	results := make(map[int]string)
	var mu sync.Mutex

	tracer.SetTraceId(customTraceId, func() {
		// 启动多个goroutine
		for i := 0; i < 5; i++ {
			wg.Add(1)
			finalI := i
			tracer.Go(func() {
				defer wg.Done()
				// 在每个goroutine中获取traceId
				gotTraceId := tracer.GetTraceId()
				mu.Lock()
				results[finalI] = gotTraceId
				mu.Unlock()
			})
		}
	})

	wg.Wait()

	// 验证所有goroutine获取到的traceId都是正确的
	for i, gotTraceId := range results {
		if gotTraceId != customTraceId {
			t.Errorf("Goroutine %d traceId = %s; want %s", i, gotTraceId, customTraceId)
		}
	}
}
