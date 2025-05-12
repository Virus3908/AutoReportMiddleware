package interfaces

import (
	"context"
	"mime/multipart"

	"github.com/jackc/pgx/v5"
	"google.golang.org/protobuf/proto"
)

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

type TxManager interface {
	WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error
}

type StorageClient interface {
	UploadFileAndGetURL(ctx context.Context, file multipart.File, originalFilename string) (string, error)
	DeleteFileByURL(ctx context.Context, fileURL string) error
}

type MessageClient interface {
	SendMessage(ctx context.Context, key string, message proto.Message) error
}
