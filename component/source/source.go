package source

import (
	"context"
	"data-backup/component/target"
)

type Source interface {
	Backup(ctx context.Context, target target.Target) error
	Restore(ctx context.Context, target target.Target) error
}

var id2SourceTable = make(map[string]Source)

func registerSource(key string, source Source) {
	id2SourceTable[key] = source
}

func GetSource(id string) Source {
	return id2SourceTable[id]
}
