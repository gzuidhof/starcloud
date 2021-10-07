package middleware

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerMiddlewareConfig struct {
	LogSuccesfulRequests bool
}

func NewLoggerMiddleware(name string, cfg LoggerMiddlewareConfig, fields ...zapcore.Field) fiber.Handler {
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

			if errHandler != nil { // I'm not sure if this `if` condition is necessary - but let's be defensive.
				if err := errHandler(c, chainErr); err != nil {
					_ = c.SendStatus(fiber.StatusInternalServerError)
				}
			}
		}
		// The callsite is useless for the logger field here
		middlewareLogger := logctx.WithOptions(zap.WithCaller(false))

		resp := c.Response()
		if resp == nil { // Should never happen, sanity check.
			zap.L().Error("Response is nil in logger (would have been panic)")
			return nil
		}

		latency := time.Since(start)
		status := resp.StatusCode()
		if chainErr != nil {
			middlewareLogger.Error(
				chainErr.Error(),
				zap.Int("status", status),
				zap.Duration("latency", latency),
			)
		} else {
			if !cfg.LogSuccesfulRequests && (status == 200 || status == 201 || status == 204 || status == 400 || status == 401 || status == 403 || status == 404) {
				return nil
			}
			middlewareLogger.Info(
				"",
				zap.Int("status", status),
				zap.Duration("latency", latency),
			)
		}
		return nil
	}
}
