package middleware

import (
	"avito/internal/pkg/errors"
	"fmt"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"math/rand"
	"time"
)

type Middleware struct {
	logger *zap.Logger
}

func NewMiddleware(logger *zap.Logger) Middleware {
	return Middleware{
		logger: logger,
	}
}

func (m *Middleware) LogMiddleware(ctx *routing.Context) error {

	start := time.Now()
	reqID := fmt.Sprintf("%016x", rand.Int())[:10]
	defer func() {
		m.logger.Info(string(ctx.Path()),
			zap.String("reqId:", reqID),
			zap.String("method", string(ctx.Method())),
			zap.String("remote_addr", ctx.RemoteAddr().String()),
			zap.Time("start", start),
			zap.Duration("work_time", time.Since(start)),
		)
	}()

	err := ctx.Next()
	if err != nil {
		m.logger.Error("Handle error",
			zap.String("reqId:", fmt.Sprintf("%v", reqID)),
			zap.String("error:", err.Error()),
		)
		return errors.ErrorHandler(ctx, err)
	}

	return nil
}

func (m *Middleware) PanicMiddleware(ctx *routing.Context) error {
	defer func() {
		if err := recover(); err != nil {
			m.logger.Error("Handle panic",
				zap.String("error: ", fmt.Sprintf("%v", err)),
			)
			ctx.Error( "Don't worry, we will be available soon!", fasthttp.StatusInternalServerError)
		}
	}()

	err := ctx.Next()
	if err != nil {
		m.logger.Error("Handle Internal error",
			zap.String("error: ", fmt.Sprintf("%v", err)),
		)
		return routing.NewHTTPError(fasthttp.StatusInternalServerError, "Don't worry, we will be available soon!")
	}
	return nil
}
