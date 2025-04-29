package logger

type LogField struct {
	Key   string
	Value interface{}
}

type Logger interface {
	Info(msg string, fields ...LogField)
	Error(msg string, fields ...LogField)
	Warn(msg string, fields ...LogField)
	Fatal(msg string, fields ...LogField)
	Debug(msg string, fields ...LogField)
	Sync()
}

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	PanicLevel = "panic"
)