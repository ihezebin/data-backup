package cmd

import (
	"context"
	"data-backup/component/email"
	"data-backup/component/source"
	"data-backup/component/target"
	"data-backup/component/task"
	"data-backup/config"
	"data-backup/cron"
	"data-backup/server"
	"os"
	"path/filepath"
	"time"

	_ "github.com/ihezebin/oneness"
	"github.com/ihezebin/oneness/logger"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var (
	configPath string
)

func Run(ctx context.Context) error {

	app := &cli.App{
		Name:    "data-backup",
		Version: "v1.0.1",
		Usage:   "Rapid construction template of Web service based on DDD architecture",
		Authors: []*cli.Author{
			{Name: "hezebin", Email: "ihezebin@qq.com"},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Destination: &configPath,
				Name:        "config", Aliases: []string{"c"},
				Value: "./config/config.toml",
				Usage: "config file path (default find file from pwd and exec dir",
			},
		},
		Before: func(c *cli.Context) error {
			if configPath == "" {
				return errors.New("config path is empty")
			}

			conf, err := config.Load(configPath)
			if err != nil {
				return errors.Wrapf(err, "load config error, path: %s", configPath)
			}
			logger.Debugf(ctx, "load config: %+v", conf.String())

			if err = initComponents(ctx, conf); err != nil {
				return errors.Wrap(err, "init components error")
			}

			if err = cron.Init(ctx, conf); err != nil {
				return errors.Wrap(err, "init cron error")
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			if err := cron.Run(ctx); err != nil {
				logger.WithError(err).Fatalf(ctx, "cron run error")
			}

			if err := server.Run(ctx, config.GetConfig().Port); err != nil {
				logger.WithError(err).Fatalf(ctx, "server run error, port: %d", config.GetConfig().Port)
			}

			return nil
		},
	}

	return app.Run(os.Args)
}

func initComponents(ctx context.Context, conf *config.Config) error {
	// init logger
	if conf.Logger != nil {
		logger.ResetLoggerWithOptions(
			logger.WithServiceName(conf.ServiceName),
			logger.WithPrettyCallerHook(),
			logger.WithTimestampHook(),
			logger.WithLevel(conf.Logger.Level),
			//logger.WithLocalFsHook(filepath.Join(conf.Pwd, conf.Logger.Filename)),
			// 每天切割，保留 3 天的日志
			logger.WithRotateLogsHook(filepath.Join(conf.Pwd, conf.Logger.Filename), time.Hour*24, time.Hour*24*3),
		)
	}

	// init email
	if conf.Email != nil {
		if err := email.Init(*conf.Email); err != nil {
			return errors.Wrap(err, "init email client error")
		}
	}

	// init sources
	if err := source.RegisterMongoSources(ctx, conf.MongoSources); err != nil {
		return errors.Wrap(err, "init mongo sources error")
	}
	if err := source.RegisterMysqlSources(ctx, conf.MysqlSources); err != nil {
		return errors.Wrap(err, "init mongo sources error")
	}

	// init targets
	if err := target.RegisterOSSTargets(ctx, conf.OSSTargets); err != nil {
		return errors.Wrap(err, "init oss targets error")
	}

	// init tasks
	if err := task.RegisterTasks(ctx, conf.DefaultCron, conf.Tasks); err != nil {
		return errors.Wrap(err, "init tasks error")
	}

	return nil
}
