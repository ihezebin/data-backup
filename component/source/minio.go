package source

import (
	"bytes"
	"context"
	"data-backup/component/target"
	"io"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
)

type MinioSource struct {
	Id       string   `json:"id" mapstructure:"id"`
	DSN      string   `json:"dsn" mapstructure:"dsn"`
	Prefixes []string `json:"prefixes" mapstructure:"prefixes"`
	Client   *minio.Client
	bucket   string
}

func RegisterMinioSources(ctx context.Context, sources []*MinioSource) error {
	for _, source := range sources {
		dsn := source.DSN
		u, err := url.Parse(dsn)
		if err != nil {
			return errors.Wrap(err, "mongo dsn parse error")
		}

		bucket := strings.Trim(u.Path, "/")
		if bucket == "" {
			return errors.New("missing bucket")
		}

		source.bucket = bucket

		accessKey := u.User.Username()
		secretKey, _ := u.User.Password()

		client, err := minio.New(u.Host, &minio.Options{
			Creds: credentials.NewStaticV4(accessKey, secretKey, ""),
		})
		if err != nil {
			return errors.Wrapf(err, "new minio client err")
		}

		exists, err := client.BucketExists(ctx, bucket)
		if err != nil {
			return errors.Wrapf(err, "minio bucket exists err")
		}
		if !exists {
			return errors.New("bucket not exists")
		}

		source.Client = client

		registerSource(source.Id, source)
	}
	return nil
}

var _ Source = (*MinioSource)(nil)

func (s *MinioSource) Backup(ctx context.Context, target target.Target) error {
	for _, prefix := range s.Prefixes {
		objsCh := s.Client.ListObjects(ctx, s.bucket, minio.ListObjectsOptions{
			Prefix:    prefix,
			Recursive: true,
		})
		for objInfo := range objsCh {
			obj, err := s.Client.GetObject(ctx, s.bucket, objInfo.Key, minio.GetObjectOptions{})
			if err != nil {
				return errors.Wrapf(err, "get object error, bucket: %s, key: %s", s.bucket, objInfo.Key)
			}

			objData, err := io.ReadAll(obj)
			if err != nil {
				return errors.Wrapf(err, "read object error, bucket: %s, key: %s", s.bucket, objInfo.Key)
			}

			importKey := path.Join(s.bucket, objInfo.Key)
			err = target.Import(ctx, importKey, objData)
			if err != nil {
				return errors.Wrapf(err, "target import error, bucket: %s, import key: %s", s.bucket, importKey)
			}
		}

	}
	return nil
}

func (s *MinioSource) Restore(ctx context.Context, target target.Target) error {
	for _, prefix := range s.Prefixes {
		exportKey := path.Join(s.bucket, prefix)
		objs, err := target.ExportMulti(ctx, exportKey)
		if err != nil {
			return errors.Wrapf(err, "export multi error, exportKey: %s", exportKey)
		}

		for _, obj := range objs {
			key, err := filepath.Rel(s.bucket, obj.Key)
			if err != nil {
				return errors.Wrapf(err, "filepath rel error, bucket: %s, key: %s", s.bucket, obj.Key)
			}

			_, err = s.Client.PutObject(ctx, s.bucket, key, bytes.NewReader(obj.Data), int64(len(obj.Data)), minio.PutObjectOptions{})
			if err != nil {
				return errors.Wrapf(err, "put object error, bucket: %s, key: %s, data length: %d", s.bucket, obj.Key, len(obj.Data))
			}
		}
	}

	return nil
}
