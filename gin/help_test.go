package gin

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestTimeStat(t *testing.T) {
	// 清理全局map，避免其他测试的影响
	mapRecord = make(map[string]*timeUseRecord)

	// 测试基本功能和计数器
	t.Run("basic functionality and counter", func(t *testing.T) {
		start := time.Now()
		time.Sleep(10 * time.Millisecond) // 模拟执行时间
		TimeStat("test1", start)

		if record := mapRecord["test1"]; record == nil {
			t.Error("Expected record to be created")
		} else if record.Count != 1 {
			t.Errorf("Expected count to be 1, got %d", record.Count)
		} else if record.UseTime < 10*time.Millisecond {
			t.Errorf("Expected use time to be at least 10ms, got %v", record.UseTime)
		}
	})

	// 测试多次调用和平均时间
	t.Run("multiple calls and average time", func(t *testing.T) {
		start := time.Now()
		time.Sleep(10 * time.Millisecond)
		TimeStat("test2", start)

		start = time.Now()
		time.Sleep(20 * time.Millisecond)
		TimeStat("test2", start)

		if record := mapRecord["test2"]; record == nil {
			t.Error("Expected record to be created")
		} else if record.Count != 2 {
			t.Errorf("Expected count to be 2, got %d", record.Count)
		} else {
			avg := record.UseTime / time.Duration(record.Count)
			if avg < 15*time.Millisecond { // 平均应该在15ms左右
				t.Errorf("Expected average time to be around 15ms, got %v", avg)
			}
		}
	})

	// 测试自定义日志函数
	t.Run("custom log function", func(t *testing.T) {
		var buf bytes.Buffer
		customLog := func(format string, args ...any) {
			fmt.Fprintf(&buf, format, args...)
		}

		start := time.Now()
		time.Sleep(10 * time.Millisecond)
		TimeStat("test3", start, customLog)

		logOutput := buf.String()
		if logOutput == "" {
			t.Error("Expected log output, got empty string")
		}
		if record := mapRecord["test3"]; record == nil {
			t.Error("Expected record to be created")
		} else if !bytes.Contains(buf.Bytes(), []byte(fmt.Sprintf("test3 executed %d times", record.Count))) {
			t.Error("Log output doesn't contain execution count")
		}
	})
}
