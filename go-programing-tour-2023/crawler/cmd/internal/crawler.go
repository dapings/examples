package internal

import (
	"errors"
	"time"

	"github.com/dapings/examples/go-programing-tour-2023/crawler/collect"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/engine"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/limiter"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/log"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/proxy"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/spider"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/sqldb"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/storage/sqlstorage"
	"github.com/go-micro/plugins/v4/config/encoder/toml"
	"go-micro.dev/v4/config"
	"go-micro.dev/v4/config/reader"
	"go-micro.dev/v4/config/reader/json"
	"go-micro.dev/v4/config/source"
	"go-micro.dev/v4/config/source/file"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/time/rate"
)

func LoadConfig() (config.Config, error) {
	// load config
	enc := toml.NewEncoder()
	cfg, cfgErr := config.NewConfig(config.WithReader(json.NewReader(reader.WithEncoder(enc))))
	if cfgErr != nil {
		return nil, cfgErr
	}

	err := cfg.Load(file.NewSource(file.WithPath("config.toml"), source.WithEncoder(enc)))
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func ConfigLogger(cfg config.Config) (*zap.Logger, error) {
	if cfg == nil {
		return nil, errors.New("config nil")
	}

	logText := cfg.Get("logLevel").String("INFO")
	logLevel, err := zapcore.ParseLevel(logText)
	if err != nil {
		return nil, err
	}

	// logger
	plugin := log.NewStdoutPlugin(logLevel)
	logger := log.NewLogger(plugin)
	logger.Info("logger init")

	// set zap global logger
	zap.ReplaceGlobals(logger)

	return logger, nil
}

func ConfigProxyFunc(cfg config.Config, logger *zap.Logger) (proxy.Func, int, error) {
	// proxy
	var err error
	var proxyFunc proxy.Func

	proxyURLs := cfg.Get("fetcher", "proxy").StringSlice([]string{})
	timeout := cfg.Get("fetcher", "timeout").Int(5000)
	logger.Sugar().Info("proxy list:", proxyURLs, " timeout: ", timeout)

	if proxyFunc, err = proxy.RoundRobinProxySwitcher(proxyURLs...); err != nil {
		logger.Error("round robin proxy switcher failed", zap.Error(err))

		return nil, 0, err
	}

	return proxyFunc, timeout, nil
}

func ConfigFetcher(proxyFunc proxy.Func, timeout int, logger *zap.Logger) spider.Fetcher {
	// fetcher
	return collect.BrowserFetch{Timeout: time.Duration(timeout) * time.Millisecond, Proxy: proxyFunc, Logger: logger}
}

func ConfigStorager(cfg config.Config, logger *zap.Logger) (spider.Storage, error) {
	// storage
	var storager spider.Storage
	var err error

	sqlURL := cfg.Get("storage", "sqlURL").String("")
	sqldb.ConnStrWithMySQL = sqlURL

	if storager, err = sqlstorage.New(
		sqlstorage.WithSQLURL(sqldb.ConnStrWithMySQL),
		sqlstorage.WithLogger(logger.Named("sqlDB")),
		sqlstorage.WithBatchCount(2),
	); err != nil {
		logger.Error("create sql storage failed", zap.Error(err))

		return nil, err
	}

	return storager, nil
}

func ParseTaskConfig(logger *zap.Logger, f spider.Fetcher, s spider.Storage, cfgs []spider.TaskConfig) []*spider.Task {
	tasks := make([]*spider.Task, 0, 1000)
	for _, cfg := range cfgs {
		t := spider.NewTask(
			spider.WithName(cfg.Name),
			spider.WithReload(cfg.Reload),
			spider.WithCookie(cfg.Cookie),
			spider.WithLogger(logger),
			spider.WithStorage(s),
		)

		if cfg.WaitTime > 0 {
			t.WaitTime = cfg.WaitTime
		}

		if cfg.MaxDepth > 0 {
			t.MaxDepth = cfg.MaxDepth
		}

		var limits []limiter.RateLimiter
		if len(cfg.Limits) > 0 {
			for _, limitCfg := range cfg.Limits {
				// speed limiter 限速器，2秒1个，60秒20个
				l := rate.NewLimiter(limiter.Per(limitCfg.EventCount, time.Duration(limitCfg.EventDur)*time.Second), 1)
				limits = append(limits, l)
			}
			multiLimiter := limiter.Multi(limits...)
			t.Limit = multiLimiter
		}

		switch cfg.Fetcher {
		case "browser":
			t.Fetcher = f
		}

		tasks = append(tasks, t)
	}

	return tasks
}

func ConfigTasks(cfg config.Config, f spider.Fetcher, storager spider.Storage, logger *zap.Logger) (*engine.Crawler, error) {
	// init tasks
	// seeds slice cap
	var seeds = make([]*spider.Task, 0, 1000)
	var taskCfg []spider.TaskConfig
	if err := cfg.Get("Tasks").Scan(&taskCfg); err != nil {
		logger.Error("init seed tasks failed", zap.Error(err))

		return nil, err
	}

	seeds = ParseTaskConfig(logger, f, storager, taskCfg)

	s := engine.NewEngine(
		engine.WithWorkCount(5),
		engine.WithFetcher(f),
		engine.WithLogger(logger),
		engine.WithSeeds(seeds),
		engine.WithScheduler(engine.NewSchedule()),
	)

	return s, nil
}

func ConfigMasterServer(cfg config.Config, logger *zap.Logger) (*ServerConfig, error) {
	var sconfig ServerConfig
	if err := cfg.Get("MasterServer").Scan(&sconfig); err != nil {
		logger.Error("get master gRPC Server config failed", zap.Error(err))

		return nil, err
	}

	logger.Sugar().Debugf("master grpc server config, %+v", sconfig)

	return &sconfig, nil
}

func ConfigWorkerServer(cfg config.Config, logger *zap.Logger) (*ServerConfig, error) {
	var sconfig ServerConfig
	if err := cfg.Get("GRPCServer").Scan(&sconfig); err != nil {
		logger.Error("get worker gRPC Server config failed", zap.Error(err))

		return nil, err
	}

	logger.Sugar().Debugf("worker grpc server config, %+v", sconfig)

	return &sconfig, nil
}
