package logger

import (
	"errors"

	"github.com/piotrkowalczuk/zapstackdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	production  = "production"
	development = "development"
	stackdriver = "stackdriver"
)

type Opts struct {
	Environment string
	Level       string
}

// Init allocates new logger based on given options.
func Init(service, version string, opts Opts) (logger *zap.Logger, err error) {
	if service == "" {
		return nil, errors.New("logger: service name is missing, logger cannot be initialized")
	}
	if version == "" {
		return nil, errors.New("logger: service version is missing, logger cannot be initialized")
	}

	var (
		zapCfg  zap.Config
		zapOpts []zap.Option
		lvl     zapcore.Level
	)
	switch opts.Environment {
	case production:
		zapCfg = zap.NewProductionConfig()
	case stackdriver:
		zapCfg = zapstackdriver.NewStackdriverConfig()
	case development:
		zapCfg = zap.NewDevelopmentConfig()
	default:
		zapCfg = zap.NewProductionConfig()
	}

	if err = lvl.Set(opts.Level); err != nil {
		return nil, err
	}
	zapCfg.Level.SetLevel(lvl)

	logger, err = zapCfg.Build(zapOpts...)
	if err != nil {
		return nil, err
	}
	logger = logger.With(zap.Object("serviceContext", &zapstackdriver.ServiceContext{
		Service: service,
		Version: version,
	}))
	logger.Info("logger has been initialized", zap.String("environment", opts.Environment))

	return logger, nil
}
