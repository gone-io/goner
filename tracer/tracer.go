package tracer

import (
	"github.com/gone-io/gone/v2"
	"github.com/google/uuid"
	"github.com/jtolds/gls"
	"sync"
)

// Tracer Log tracking, which is used to assign a unified traceId to the same call link to facilitate log tracking.
type Tracer interface {

	//SetTraceId to set `traceId` to the calling function. If traceId is an empty string, an automatic one will
	//be generated. TraceId can be obtained by using the GetTraceId () method in the calling function.
	SetTraceId(traceId string, fn func())

	//GetTraceId Get the traceId of the current goroutine
	GetTraceId() string

	//Go Start a new goroutine instead of `go func`, which can pass the traceId to the new goroutine.
	Go(fn func())
}

var load = gone.OnceLoad(func(loader gone.Loader) error {
	return loader.Load(
		&tracer{},
		gone.IsDefault(new(Tracer)),
		gone.LazyFill(),
	)
})

func Load(loader gone.Loader) error {
	return load(loader)
}

// Priest Deprecated, use Load instead
func Priest(loader gone.Loader) error {
	return Load(loader)
}

type tracer struct {
	gone.Flag
	gone.Logger `gone:"*"`
}

var xMap sync.Map

func (t *tracer) GonerName() string {
	return "gone-tracer"
}

func (t *tracer) SetTraceId(traceId string, cb func()) {
	SetTraceId(traceId, cb, t.Warnf)
}

func (t *tracer) GetTraceId() (traceId string) {
	return GetTraceId()
}

func (t *tracer) Go(cb func()) {
	traceId := t.GetTraceId()
	if traceId == "" {
		go cb()
	} else {
		go func() {
			t.SetTraceId(traceId, cb)
		}()
	}
}

func GetTraceId() (traceId string) {
	gls.EnsureGoroutineId(func(gid uint) {
		if v, ok := xMap.Load(gid); ok {
			traceId = v.(string)
		}
	})
	return
}

func SetTraceId(traceId string, cb func(), log ...func(format string, args ...any)) {
	id := GetTraceId()
	if "" != id {
		if len(log) > 0 {
			log[0]("SetTraceId not success for Having been set")
		}
		cb()
		return
	} else {
		if traceId == "" {
			traceId = uuid.New().String()
		}
		gls.EnsureGoroutineId(func(gid uint) {
			xMap.Store(gid, traceId)
			defer xMap.Delete(gid)
			cb()
		})
	}
}

//func GetGoroutineId() (gid uint64) {
//	var (
//		buf [64]byte
//		n   = runtime.Stack(buf[:], false)
//		stk = strings.TrimPrefix(string(buf[:n]), "goroutine ")
//	)
//	idField := strings.Fields(stk)[0]
//	var err error
//	gid, err = strconv.ParseUint(idField, 10, 64)
//	if err != nil {
//		panic(fmt.Errorf("can not get goroutine id: %v", err))
//	}
//	return
//}
