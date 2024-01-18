package main

import (
	"github.com/labstack/gommon/log"
	"github.com/opentracing/opentracing-go"
	"github.com/sefikcan/ms-grpc-sample/product/internal/server"
	"github.com/sefikcan/ms-grpc-sample/product/pkg/config"
	"github.com/sefikcan/ms-grpc-sample/product/pkg/logger"
	"github.com/sefikcan/ms-grpc-sample/product/pkg/metric"
	"github.com/sefikcan/ms-grpc-sample/product/pkg/storage/mongo"
	"github.com/uber/jaeger-client-go"
	jaegerCfg "github.com/uber/jaeger-client-go/config"
	jaegerLog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"io"
)

func main() {
	cfg := config.NewConfig()

	zapLogger := logger.NewLogger(cfg)
	zapLogger.InitLogger()
	zapLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %v", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode, false)

	log.Info("Starting product api server")

	mongoDb, err := mongo.NewMongo(cfg)

	jaegerConfigInstance := jaegerCfg.Configuration{
		ServiceName: cfg.Metric.ServiceName,
		Sampler: &jaegerCfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegerCfg.ReporterConfig{
			LogSpans:           cfg.Jaeger.LogSpans,
			LocalAgentHostPort: cfg.Jaeger.Host,
		},
	}

	tracer, closer, err := jaegerConfigInstance.NewTracer(
		jaegerCfg.Logger(jaegerLog.StdLogger),
		jaegerCfg.Metrics(metrics.NullFactory),
	)
	if err != nil {
		log.Fatal("cannot create tracer", err)
	}
	zapLogger.Info("Jaeger connected")

	opentracing.SetGlobalTracer(tracer)
	defer func(closer io.Closer) {
		err := closer.Close()
		if err != nil {

		}
	}(closer)
	zapLogger.Info("Opentracing connected")

	_, err = metric.CreateMetrics(cfg.Metric.Url, cfg.Metric.ServiceName)
	if err != nil {
		zapLogger.Errorf("CreateMetrics error: %s", err)
	}

	zapLogger.Infof("Metrics available URL: %s, ServiceName: %s", cfg.Metric.Url, cfg.Metric.ServiceName)

	s := server.NewServer(cfg, mongoDb, zapLogger)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
