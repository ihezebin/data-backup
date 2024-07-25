package config

import (
	"data-backup/component/source"
	"data-backup/component/target"
	"data-backup/component/task"
	"encoding/json"
	"os"

	"github.com/ihezebin/oneness/config"
	"github.com/ihezebin/oneness/email"
	"github.com/ihezebin/oneness/logger"
	"github.com/pkg/errors"
)

type Config struct {
	ServiceName  string                `json:"service_name" mapstructure:"service_name"`
	Port         uint                  `json:"port" mapstructure:"port"`
	DefaultCron  string                `json:"default_cron" mapstructure:"default_cron"`
	Logger       *LoggerConfig         `json:"logger" mapstructure:"logger"`
	Email        *email.Config         `json:"email" mapstructure:"email"`
	Pwd          string                `json:"-" mapstructure:"-"`
	MongoSources []*source.MongoSource `json:"mongo_sources" mapstructure:"mongo_sources"`
	MysqlSources []*source.MysqlSource `json:"mysql_sources" mapstructure:"mysql_sources"`
	OSSTargets   []*target.OSSTarget   `json:"oss_targets" mapstructure:"oss_targets"`
	Tasks        []*task.Task          `json:"tasks" mapstructure:"tasks"`
}

type LoggerConfig struct {
	Level    logger.Level `json:"level" mapstructure:"level"`
	Filename string       `json:"filename" mapstructure:"filename"`
}

var gConfig *Config = &Config{}

func (c *Config) String() string {
	data, _ := json.Marshal(c)
	return string(data)
}

func GetConfig() *Config {
	return gConfig
}

func Load(path string) (*Config, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap(err, "get pwd error")
	}

	if err = config.NewWithFilePath(path).Load(gConfig); err != nil {
		return nil, errors.Wrap(err, "load config error")
	}

	gConfig.Pwd = pwd

	return gConfig, nil
}
