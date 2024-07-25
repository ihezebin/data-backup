package target

import (
	"bytes"
	"context"
	"io"
	"path"

	"github.com/ihezebin/oneness/oss"
	"github.com/pkg/errors"
)

type OSSTarget struct {
	Key    string `json:"key" mapstructure:"key"`
	Dsn    string `json:"dsn" mapstructure:"dsn"`
	Dir    string `json:"dir" mapstructure:"dir"`
	Client oss.Client
}

var key2OSSTargetTable = make(map[string]*OSSTarget)

func GetOSSTarget(key string) *OSSTarget {
	return key2OSSTargetTable[key]
}

func InitOSSTargets(ctx context.Context, targets []*OSSTarget) error {
	for _, target := range targets {
		client, err := oss.NewClient(target.Dsn)
		if err != nil {
			return errors.Wrap(err, "create oss client error")
		}
		key2OSSTargetTable[target.Key] = &OSSTarget{
			Dsn:    target.Dsn,
			Dir:    target.Dir,
			Client: client,
		}
	}
	return nil
}

var _ Target = (*OSSTarget)(nil)

func (t *OSSTarget) Store(ctx context.Context, key string, data []byte) error {
	name := path.Join(t.Dir, key)
	err := t.Client.PutObject(ctx, name, bytes.NewBuffer(data))
	if err != nil {
		return errors.Wrapf(err, "put object error, key: %s", key)
	}

	return nil
}

func (t *OSSTarget) Restore(ctx context.Context, key string) ([]byte, error) {
	name := path.Join(t.Dir, key)
	object, err := t.Client.GetObject(ctx, name)
	if err != nil {
		return nil, errors.Wrapf(err, "get object error, key: %s", key)
	}
	data, err := io.ReadAll(object)
	if err != nil {
		return nil, errors.Wrapf(err, "read object error, key: %s", key)
	}

	return data, nil
}
