package source

import (
	"context"
)

const (
	TypeMongoDB = "mongo"
	TypeMysql   = "mysql"
)

var SupportSourceTypes = map[string]struct{}{
	TypeMongoDB: {},
	TypeMysql:   {},
}

type Data struct {
	Key     string
	Content []byte
}

type Source interface {
	Export(ctx context.Context) ([]Data, error)
	Import(ctx context.Context, data Data) error
	Keys() []string
}
