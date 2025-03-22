package g

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
