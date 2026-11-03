package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log  *zap.Logger
	once sync.Once
)

// Init initializes the global zap JSON logger.
func Init(appEnv string) {
	once.Do(func() {
		level := zapcore.InfoLevel
		if appEnv == "local" || appEnv == "dev" || appEnv == "development" {
			level = zapcore.DebugLevel
		}

		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.MillisDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			level,
		)

		log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	})
}

// Get returns the singleton zap logger instance.
func Get() *zap.Logger {
	if log == nil {
		Init("local")
	}
	return log
}

// Sync flushes buffered log entries.
func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}

// Info logs a message at InfoLevel.
func Info(msg string, fields ...zap.Field) {
	Get().Info(msg, fields...)
}

// Error logs a message at ErrorLevel.
func Error(msg string, fields ...zap.Field) {
	Get().Error(msg, fields...)
}

// Debug logs a message at DebugLevel.
func Debug(msg string, fields ...zap.Field) {
	Get().Debug(msg, fields...)
}

// Warn logs a message at WarnLevel.
func Warn(msg string, fields ...zap.Field) {
	Get().Warn(msg, fields...)
}

// Fatal logs a message at FatalLevel.
func Fatal(msg string, fields ...zap.Field) {
	Get().Fatal(msg, fields...)
}
