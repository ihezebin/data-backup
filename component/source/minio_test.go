package source

import (
	"context"
	"data-backup/component/target"
	"testing"
)

func TestMinioBackup(t *testing.T) {
	ctx := context.Background()
	err := RegisterMinioSources(ctx, []*MinioSource{
		{
			Id:       "1",
			DSN:      "minio://xxx:xxx@127.0.0.1:9000/blog-minio",
			Prefixes: []string{"preview"},
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	err = target.RegisterOSSTargets(ctx, []*target.OSSTarget{
		{
			Id:  "1",
			Dsn: "cos://xxx:xxx@cos.ap-chengdu.myqcloud.com/hezebin-1258606727",
			Dir: "backup_test",
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	_source := GetSource("1")
	if _source == nil {
		t.Fatal("minio source not found")
	}

	_target := target.GetTarget("1")
	if _target == nil {
		t.Fatal("minio target not found")
	}

	err = _source.Backup(ctx, _target)

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("minio backup success")
}

func TestMinioRestore(t *testing.T) {
	ctx := context.Background()
	err := RegisterMinioSources(ctx, []*MinioSource{
		{
			Id:       "1",
			DSN:      "minio://xxxx:xxxx@127.0.0.1:9000/blog-minio",
			Prefixes: []string{"preview"},
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	err = target.RegisterOSSTargets(ctx, []*target.OSSTarget{
		{
			Id:  "1",
			Dsn: "cos://xxxx:xxxxx@cos.ap-chengdu.myqcloud.com/hezebin-1258606727",
			Dir: "backup_test",
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	_source := GetSource("1")
	if _source == nil {
		t.Fatal("minio source not found")
	}

	_target := target.GetTarget("1")
	if _target == nil {
		t.Fatal("minio target not found")
	}

	err = _source.Restore(ctx, _target)

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("minio restore success")
}
