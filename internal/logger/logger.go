package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
)

type ZapLogger struct {
	logger *zap.Logger
}

func NewLogger(logLevel string) (*ZapLogger, error) {
	var level zap.AtomicLevel

	switch logLevel {
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
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	config := zap.Config{
		Level:            level,
		Development:      false,
		Sampling:         nil,
		Encoding:         "json",
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

func (l *ZapLogger) Info(msg string, fields ...LogField) {
	l.logger.Info(msg, l.toZapFields(fields)...)
}
func (l *ZapLogger) Error(msg string, fields ...LogField) {
	l.logger.Error(msg, l.toZapFields(fields)...)
}
func (l *ZapLogger) Warn(msg string, fields ...LogField) {
	l.logger.Warn(msg, l.toZapFields(fields)...)
}
func (l *ZapLogger) Fatal(msg string, fields ...LogField) {
	l.logger.Fatal(msg, l.toZapFields(fields)...)
}
func (l *ZapLogger) Debug(msg string, fields ...LogField) {
	l.logger.Debug(msg, l.toZapFields(fields)...)
}

func (l *ZapLogger) toZapFields(fields []LogField) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = zap.Any(f.Key, f.Value)
	}
	return zapFields
}

func (l *ZapLogger) LoggingMidleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// l.Info("Received request",
		// 	LogField{Key: "method", Value: r.Method},
		// 	LogField{Key: "path", Value: r.URL.Path},
		// )
		l.logger.Info("Received request",
			zap.String("method", r.Method), zap.String("path", r.URL.Path))
		next.ServeHTTP(w, r)
	})
}
