package target

import (
	"context"
)

type Target interface {
	Import(ctx context.Context, key string, data []byte) error
	Export(ctx context.Context, key string) ([]byte, error)
}

var key2TargetTable = make(map[string]Target)

func registerTarget(key string, target Target) {
	key2TargetTable[key] = target
}

func GetTarget(key string) Target {
	return key2TargetTable[key]
}
