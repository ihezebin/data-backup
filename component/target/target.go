package target

import (
	"context"
)

const (
	TypeOss = "oss"
)

var SupportTargetTypes = map[string]struct{}{
	TypeOss: {},
}

type Target interface {
	Store(ctx context.Context, key string, data []byte) error
	Restore(ctx context.Context, sourceKey string) ([]byte, error)
}
