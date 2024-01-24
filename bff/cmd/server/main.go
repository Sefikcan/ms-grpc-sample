package main

import (
	"github.com/sefikcan/ms-grpc-sample/bff/internal"
	"github.com/sefikcan/ms-grpc-sample/bff/pkg/config"
	"github.com/sefikcan/ms-grpc-sample/bff/pkg/logger"
	_ "github.com/sefikcan/ms-grpc-sample/docs"
	"log"
)

// @title Ms gRPC Sample
// @version 1.0
// @description Ms gRPC Sample
// @contact.name Sefik Can Kanber
// @contact.url https://github.com/sefikcan
// @BasePath /api/v1
// @host localhost:50050
func main() {
	log.Println("Starting bff api server")

	cfg := config.NewConfig()

	zapLogger := logger.NewLogger(cfg)
	zapLogger.InitLogger()
	zapLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %v", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode, false)

	s := internal.NewServer(cfg, zapLogger)
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
