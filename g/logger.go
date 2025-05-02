package g

import (
	"context"
	"github.com/gone-io/gone/v2"
)

// CtxLogger Log tracking, which is used to assign a unified traceId to the same call link to facilitate log tracking.
// Examples:
//
//	type user struct {
//		gone.Flag
//		logger CtxLogger `gone:"*"` //Inject  Logger
//	}
//
//	func (u *user) Use(ctx context.Context) (err error) {
//		// get traceId from openTelemetry context and inject it into the logger
//		logger := u.logger.Ctx(ctx)
//
//		logger.Infof("hello")
//
//		return
//	}
type CtxLogger interface {
	Ctx(ctx context.Context) gone.Logger
}
