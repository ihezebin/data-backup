package target

import (
	"bytes"
	"context"
	"io"
	"path"
	"path/filepath"

	"github.com/ihezebin/oneness/oss"
	"github.com/pkg/errors"
)

type OSSTarget struct {
	Id     string `json:"id" mapstructure:"id"`
	Dsn    string `json:"dsn" mapstructure:"dsn"`
	Dir    string `json:"dir" mapstructure:"dir"`
	Client oss.Client
}

func RegisterOSSTargets(_ context.Context, targets []*OSSTarget) error {
	for _, target := range targets {
		client, err := oss.NewClient(target.Dsn)
		if err != nil {
			return errors.Wrap(err, "create oss client error")
		}

		target.Client = client
		registerTarget(target.Id, target)
	}
	return nil
}

var _ Target = (*OSSTarget)(nil)

func (t *OSSTarget) Import(ctx context.Context, key string, data []byte) error {
	name := path.Join(t.Dir, key)
	err := t.Client.PutObject(ctx, name, bytes.NewBuffer(data))
	if err != nil {
		return errors.Wrapf(err, "put object error, key: %s", key)
	}

	return nil
}

func (t *OSSTarget) ExportMulti(ctx context.Context, prefix string) ([]Data, error) {
	name := path.Join(t.Dir, prefix)
	objects, err := t.Client.GetObjects(ctx, name)
	if err != nil {
		return nil, errors.Wrapf(err, "get objects error, prefix: %s", prefix)
	}

	data := make([]Data, 0)
	for _, object := range objects {
		objData, err := io.ReadAll(object.Data)
		if err != nil {
			return nil, errors.Wrapf(err, "read object error, prefix: %s, object: %s", prefix, object)
		}
		key, err := filepath.Rel(t.Dir, object.Key)
		if err != nil {
			return nil, errors.Wrapf(err, "filepath rel error, prefix: %s, key: %s", prefix, object.Key)
		}
		data = append(data, Data{Key: key, Data: objData})
	}

	return data, nil
}

func (t *OSSTarget) Export(ctx context.Context, key string) ([]byte, error) {
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
