package logger

import (
	"context"
	"log"
	"main/internal/common/interfaces"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	PanicLevel = "panic"
)

type LoggerConfig struct {
	LogLevel string `yaml:"log_level"`
	Encoding string `yaml:"encoding"`
}

type ZapLogger struct {
	logger *zap.Logger
}

func NewLogger(cfg LoggerConfig) (*ZapLogger, error) {
	var level zap.AtomicLevel

	switch cfg.LogLevel {
	case DebugLevel:
		level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case InfoLevel:
		level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case WarnLevel:
		level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case ErrorLevel:
		level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case PanicLevel:
		level = zap.NewAtomicLevelAt(zapcore.PanicLevel)
	default:
		level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "trace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	config := zap.Config{
		Level:            level,
		Development:      false,
		Sampling:         nil,
		Encoding:         cfg.Encoding,
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		log.Printf("failed build zap log: %v", err)
		return nil, err
	}

	zap.RedirectStdLog(logger)

	return &ZapLogger{
		logger: logger,
	}, nil
}

func (l *ZapLogger) Sync() {
	if l.logger != nil {
		_ = l.logger.Sync()
	}
}

func (l *ZapLogger) Info(msg string, fields ...interfaces.LogField) {
	l.logger.Info(msg, l.toZapFields(fields)...)
}
func (l *ZapLogger) Error(msg string, fields ...interfaces.LogField) {
	l.logger.Error(msg, l.toZapFields(fields)...)
}
func (l *ZapLogger) Warn(msg string, fields ...interfaces.LogField) {
	l.logger.Warn(msg, l.toZapFields(fields)...)
}
func (l *ZapLogger) Fatal(msg string, fields ...interfaces.LogField) {
	l.logger.Fatal(msg, l.toZapFields(fields)...)
}
func (l *ZapLogger) Debug(msg string, fields ...interfaces.LogField) {
	l.logger.Debug(msg, l.toZapFields(fields)...)
}

func (l *ZapLogger) toZapFields(fields []interfaces.LogField) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = zap.Any(f.Key, f.Value)
	}
	return zapFields
}


// delete later

func (l *ZapLogger) LoggingMidleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// l.Info("Received request",
		// 	interfaces.LogField{Key: "method", Value: r.Method},
		// 	interfaces.LogField{Key: "path", Value: r.URL.Path},
		// )
		l.logger.Info("Received request",
			zap.String("method", r.Method), zap.String("path", r.URL.Path))
		next.ServeHTTP(w, r)
	})
}

func ContextWithLogger(baseLogger interfaces.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "logger", baseLogger)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetLoggerFromContext(ctx context.Context) interfaces.Logger {
	logger, ok := ctx.Value("logger").(interfaces.Logger)
	if !ok {
		return nil
	}
	return logger
}