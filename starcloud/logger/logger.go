package logger

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func SetupLogger() {
	var z zap.Config

	if viper.GetBool("development") {
		z = zap.NewDevelopmentConfig()
	} else {
		z = zap.NewProductionConfig()
		z.DisableStacktrace = true
		z.EncoderConfig.EncodeTime = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.UTC().Format(time.RFC3339)) // Let's always have UTC timestamps
		})
	}
	zl, err := z.Build()
	if err != nil {
		log.Fatalf("Failed to create global zap logger: %v", err)
	}

	zap.ReplaceGlobals(zl)
}

func GetLogger(ctx *fiber.Ctx) *zap.Logger {
	return ctx.Locals("logger").(*zap.Logger)
}
func GetSugarLogger(ctx *fiber.Ctx) *zap.SugaredLogger {
	return ctx.Locals("logger").(*zap.Logger).Sugar()
}
