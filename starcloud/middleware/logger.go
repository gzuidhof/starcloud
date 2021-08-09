package middleware

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLoggerMiddleware(name string, fields ...zapcore.Field) fiber.Handler {
	var errHandler fiber.ErrorHandler
	var zaplogger = zap.L()
	var once sync.Once
	defer zaplogger.Sync()

	if name != "" {
		zaplogger = zaplogger.Named(name)
	}
	zaplogger = zaplogger.With(fields...)

	// Return new handler
	return func(c *fiber.Ctx) (err error) {
		defer zaplogger.Sync()
		start := time.Now().UTC()

		logctx := zaplogger.With(
			// We assume a request id has been set
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.String("request_id", c.Locals("request_id").(string)),
		)

		c.Locals("logger", logctx)

		// Handle request, store err for logging
		chainErr := c.Next()
		// Manually call error handler
		if chainErr != nil {
			once.Do(func() {
				// override error handler
				errHandler = c.App().Config().ErrorHandler
			})
				

			if err := errHandler(c, chainErr); err != nil {
				_ = c.SendStatus(fiber.StatusInternalServerError)
			}
		}
		// The callsite is useless for the logger field here
		middlewareLogger := logctx.WithOptions(zap.WithCaller(false))

		stop := time.Now()
		if chainErr != nil {
			middlewareLogger.Error(
				chainErr.Error(),
				zap.Int("status", c.Response().StatusCode()),
				zap.Duration("latency", stop.Sub(start)),
			)
		} else {
			middlewareLogger.Info(
				"",
				zap.Int("status", c.Response().StatusCode()),
				zap.Duration("latency", stop.Sub(start)),
			)
		}
		return nil
	}
}