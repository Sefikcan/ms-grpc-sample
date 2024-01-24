package main

import (
	"context"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/sefikcan/ms-grpc-sample/product/internal/repository"
	"github.com/sefikcan/ms-grpc-sample/product/internal/use_case"
	"github.com/sefikcan/ms-grpc-sample/product/pkg/config"
	"github.com/sefikcan/ms-grpc-sample/product/pkg/logger"
	"github.com/sefikcan/ms-grpc-sample/product/pkg/storage/mongo"
	"google.golang.org/grpc"
	"net"
)

func main() {
	log.Info("Starting product api server")

	cfg := config.NewConfig()

	zapLogger := logger.NewLogger(cfg)
	zapLogger.InitLogger()
	zapLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %v", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode, false)

	serverAddress := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)

	listen, err := net.Listen(cfg.Server.NetworkType, serverAddress)
	if err != nil {
		zapLogger.Fatalf("Failed to listen on: %v\n", err)
	}
	zapLogger.Infof("Listening on %s\n", serverAddress)

	grpcServer := grpc.NewServer()

	db, err := mongo.NewMongo(cfg)

	databases, err := db.ListDatabaseNames(context.Background(), nil)
	if err != nil {
		zapLogger.Error("ERROR:", err.Error())
	}

	dbExists := false
	for _, dbName := range databases {
		if dbName == "productDb" {
			dbExists = true
			break
		}
	}

	// If the database doesn't exist, create it
	if !dbExists {
		err := db.Database("productDb").CreateCollection(context.Background(), "product")
		if err != nil {
			zapLogger.Error("ERROR:", err.Error())
		}
	}

	collections, err := db.Database("productDb").ListCollectionNames(context.Background(), nil)
	if err != nil {
		zapLogger.Error("ERROR:", err.Error())
	}

	collectionExists := false
	for _, colName := range collections {
		if colName == "product" {
			collectionExists = true
			break
		}
	}

	// If the collection doesn't exist, create it
	if !collectionExists {
		err := db.Database("productDb").CreateCollection(context.Background(), "product")
		if err != nil {
			zapLogger.Error("ERROR:", err.Error())
		}
	}

	productRepository := repository.NewProductRepository(db)

	use_case.NewProductUseCase(cfg, productRepository, zapLogger, grpcServer)

	zapLogger.Infof("Server started at %v", listen.Addr().String())

	err = grpcServer.Serve(listen)
	if err != nil {
		zapLogger.Error("ERROR:", err.Error())
	}
}
