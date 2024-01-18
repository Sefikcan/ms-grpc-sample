package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sefikcan/ms-grpc-sample/product/pkg/config"
	"github.com/sefikcan/ms-grpc-sample/product/pkg/logger"
	pb "github.com/sefikcan/ms-grpc-sample/proto"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	echo    *echo.Echo
	cfg     *config.Config
	mongoDB *mongo.Client
	logger  logger.Logger
}

func (s *Server) Run() error {
	log.Println("Product Service started")

	serverAddress := fmt.Sprintf("%s:%s", s.cfg.Server.Host, s.cfg.Server.Port)

	listen, err := net.Listen(s.cfg.Server.NetworkType, serverAddress)
	if err != nil {
		log.Println("ERROR:", err.Error())
	}

	grpcServer := grpc.NewServer()
	pb.RegisterProductServiceServer(grpcServer, pb.UnimplementedProductServiceServer{})

	log.Printf("Server started at %v", listen.Addr().String())

	err = grpcServer.Serve(listen)
	if err != nil {
		log.Println("ERROR:", err.Error())
	}

	return nil
}

func NewServer(cfg *config.Config, mongoDB *mongo.Client, logger logger.Logger) *Server {
	return &Server{
		echo:    echo.New(),
		cfg:     cfg,
		mongoDB: mongoDB,
		logger:  logger,
	}
}
