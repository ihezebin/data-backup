package source

import (
	"context"
	"data-backup/component/storage"
	"data-backup/component/target"
	"path"
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlSource struct {
	Id     string   `json:"id" mapstructure:"id"`
	DSN    string   `json:"dsn" mapstructure:"dsn"`
	Tables []string `json:"tables" mapstructure:"tables"`
	DB     *gorm.DB
	DBName string
}

var _ Source = (*MysqlSource)(nil)

func (s *MysqlSource) Backup(ctx context.Context, target2 target.Target) error {
	return errors.New("mysql export not supported")
}

func (s *MysqlSource) Restore(ctx context.Context, target target.Target) error {
	return errors.New("mysql import not supported")
}

func (s *MysqlSource) Keys() []string {
	keys := make([]string, 0)

	for _, table := range s.Tables {
		keys = append(keys, path.Join(s.DB.Config.Name(), table))
	}
	return keys
}

func RegisterMysqlSources(_ context.Context, sources []*MysqlSource) error {
	for _, source := range sources {
		dsn := source.DSN
		mysqlDsn, err := storage.ParseMysqlDSN(dsn)
		if err != nil {
			return errors.Wrap(err, "mysql parse dsn error")
		}

		db, err := gorm.Open(mysql.Open(mysqlDsn.FormatDSN()), &gorm.Config{})
		if err != nil {
			return errors.Wrap(err, "mysql connect error")
		}

		// https://gorm.io/docs/generic_interface.html#Connection-Pool
		sqlDB, err := db.DB()
		if err != nil {
			return errors.Wrap(err, "mysql get sql db error")
		}

		sqlDB.SetConnMaxIdleTime(time.Minute * 1)
		sqlDB.SetMaxIdleConns(0)
		sqlDB.SetMaxOpenConns(50)

		source.DB = db

		source.DBName = mysqlDsn.DBName

		registerSource(source.Id, source)
	}
	return nil
}
