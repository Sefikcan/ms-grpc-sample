package logger

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/sefikcan/ms-grpc-sample/product/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type Logger interface {
	InitLogger()
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	DPanic(args ...interface{})
	Fatal(args ...interface{})

	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	DPanicf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
}

type logger struct {
	cfg           *config.Config
	sugarLogger   *zap.SugaredLogger
	elasticClient *elastic.Client
}

var loggerLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func (l *logger) getLogLevel(cfg *config.Config) zapcore.Level {
	level, exist := loggerLevelMap[cfg.Logger.Level]
	if !exist {
		return zapcore.DebugLevel
	}

	return level
}

func (l *logger) InitLogger() {
	logLevel := l.getLogLevel(l.cfg)
	logWriter := zapcore.AddSync(os.Stderr)

	var encoderCfg zapcore.EncoderConfig
	if l.cfg.Server.Mode == "Dev" {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}

	var encoder zapcore.Encoder
	encoderCfg.LevelKey = "LEVEL"
	encoderCfg.CallerKey = "CALLER"
	encoderCfg.TimeKey = "TIME"
	encoderCfg.NameKey = "NAME"
	encoderCfg.MessageKey = "MESSAGE"

	if l.cfg.Logger.Encoding == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(encoder, logWriter, zap.NewAtomicLevelAt(logLevel))
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))

	l.sugarLogger = logger.Sugar()
	if err := l.sugarLogger.Sync(); err != nil {
		l.sugarLogger.Error(err)
	}

	elasticClient, err := elastic.NewClient(
		elastic.SetURL(l.cfg.Logger.ElasticSearchUrl),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetSniff(false),
		elastic.SetHealthcheckTimeout(5*time.Second))
	if err != nil {
		l.sugarLogger.Fatal("Failed to initialize Elasticsearch client: ", err)
	}
	l.elasticClient = elasticClient
}

func (l *logger) sendLogToElasticSearch(level zapcore.Level, message string) {
	logEntry := map[string]interface{}{
		"level":   level.String(),
		"message": message,
	}

	_, err := l.elasticClient.Index().Index("product_log_index").BodyJson(logEntry).Do(context.Background())
	if err != nil {
		l.sugarLogger.Error("Failed to send log entry to Elasticsearch: ", err)
	}
}

func (l logger) Debug(args ...interface{}) {
	message := fmt.Sprint(args...)
	l.sugarLogger.Debug(args)
	l.sendLogToElasticSearch(zapcore.DebugLevel, message)
}

func (l logger) Info(args ...interface{}) {
	message := fmt.Sprint(args...)
	l.sugarLogger.Info(args)
	l.sendLogToElasticSearch(zapcore.InfoLevel, message)
}

func (l logger) Warn(args ...interface{}) {
	message := fmt.Sprint(args...)
	l.sugarLogger.Warn(args)
	l.sendLogToElasticSearch(zapcore.WarnLevel, message)
}

func (l logger) Error(args ...interface{}) {
	message := fmt.Sprint(args...)
	l.sugarLogger.Error(args)
	l.sendLogToElasticSearch(zapcore.ErrorLevel, message)
}

func (l logger) DPanic(args ...interface{}) {
	message := fmt.Sprint(args...)
	l.sugarLogger.DPanic(args)
	l.sendLogToElasticSearch(zapcore.DPanicLevel, message)
}

func (l logger) Fatal(args ...interface{}) {
	message := fmt.Sprint(args...)
	l.sugarLogger.Fatal(args)
	l.sendLogToElasticSearch(zapcore.FatalLevel, message)
}

func (l logger) Debugf(template string, args ...interface{}) {
	message := fmt.Sprint(args...)
	l.sugarLogger.Debugf(template, args...)
	l.sendLogToElasticSearch(zapcore.DebugLevel, message)
}

func (l logger) Infof(template string, args ...interface{}) {
	message := fmt.Sprint(args...)
	l.sugarLogger.Infof(template, args...)
	l.sendLogToElasticSearch(zapcore.InfoLevel, message)
}

func (l logger) Warnf(template string, args ...interface{}) {
	message := fmt.Sprint(args...)
	l.sugarLogger.Warnf(template, args...)
	l.sendLogToElasticSearch(zapcore.WarnLevel, message)
}

func (l logger) Errorf(template string, args ...interface{}) {
	message := fmt.Sprint(args...)
	l.sugarLogger.Errorf(template, args...)
	l.sendLogToElasticSearch(zapcore.ErrorLevel, message)
}

func (l logger) DPanicf(template string, args ...interface{}) {
	message := fmt.Sprint(args...)
	l.sugarLogger.DPanicf(template, args...)
	l.sendLogToElasticSearch(zapcore.DPanicLevel, message)
}

func (l logger) Fatalf(template string, args ...interface{}) {
	message := fmt.Sprint(args...)
	l.sugarLogger.Fatalf(template, args...)
	l.sendLogToElasticSearch(zapcore.FatalLevel, message)
}

func NewLogger(cfg *config.Config) Logger {
	return &logger{
		cfg: cfg,
	}
}
