package target

import (
	"context"
)

type Target interface {
	Import(ctx context.Context, key string, data []byte) error
	Export(ctx context.Context, key string) ([]byte, error)
	ExportMulti(ctx context.Context, prefix string) ([]Data, error)
}
type Data struct {
	Key  string
	Data []byte
}

var id2TargetTable = make(map[string]Target)

func registerTarget(key string, target Target) {
	id2TargetTable[key] = target
}

func GetTarget(id string) Target {
	return id2TargetTable[id]
}
