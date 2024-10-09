package source

import (
	"context"
	"data-backup/component/target"
)

type Source interface {
	Backup(ctx context.Context, target target.Target) error
	Restore(ctx context.Context, target target.Target) error
}

var key2SourceTable = make(map[string]Source)

func registerSource(key string, source Source) {
	key2SourceTable[key] = source
}

func GetSource(key string) Source {
	return key2SourceTable[key]
}
