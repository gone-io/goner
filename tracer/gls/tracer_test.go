package gls

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestTracer_SetAndGetTraceId(t *testing.T) {
	tracer := &tracer{}

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

func TestTracer_Go(t *testing.T) {
	tracer := &tracer{}
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

func TestTracer_MultipleGoroutines(t *testing.T) {
	tracer := &tracer{}
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

func TestTracer_NestedSetTraceId(t *testing.T) {
	tracer := &tracer{}
	outerTraceId := "outer-trace-id"
	innerTraceId := "inner-trace-id"

	var outerGotTraceId, innerGotTraceId string

	tracer.SetTraceId(outerTraceId, func() {
		outerGotTraceId = tracer.GetTraceId()

		// 在已有traceId的情况下再次设置traceId
		tracer.SetTraceId(innerTraceId, func() {
			innerGotTraceId = tracer.GetTraceId()
		})
	})

	// 验证内部和外部的traceId是否正确
	if outerGotTraceId != outerTraceId {
		t.Errorf("Outer traceId = %s; want %s", outerGotTraceId, outerTraceId)
	}

	if innerGotTraceId != innerTraceId {
		t.Errorf("Inner traceId = %s; want %s", innerGotTraceId, innerTraceId)
	}
}

func TestLoad(t *testing.T) {
	gone.
		NewApp(Load).
		Run(func(tracer g.Tracer) {
			assert.NotNil(t, tracer)
		})
}
