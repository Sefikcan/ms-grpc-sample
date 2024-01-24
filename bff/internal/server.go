package internal

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sefikcan/ms-grpc-sample/bff/internal/middlewares"
	"github.com/sefikcan/ms-grpc-sample/bff/internal/product/handlers"
	"github.com/sefikcan/ms-grpc-sample/bff/pkg/config"
	"github.com/sefikcan/ms-grpc-sample/bff/pkg/logger"
	"github.com/sefikcan/ms-grpc-sample/bff/pkg/util"
	pb "github.com/sefikcan/ms-grpc-sample/proto"
	echoSwagger "github.com/swaggo/echo-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	echo   *echo.Echo
	cfg    *config.Config
	logger logger.Logger
}

func (s *Server) Run() error {
	server := &http.Server{
		Addr:           fmt.Sprintf("%s:%s", s.cfg.Server.Host, s.cfg.Server.Port),
		ReadTimeout:    time.Second * time.Duration(s.cfg.Server.ReadTimeout),
		WriteTimeout:   time.Second * time.Duration(s.cfg.Server.WriteTimeout),
		MaxHeaderBytes: s.cfg.Server.MaxHeaderBytes,
	}

	go func() {
		s.logger.Infof("Server is listening on PORT: %s", s.cfg.Server.Port)
		if err := s.echo.StartServer(server); err != nil {
			s.logger.Fatalf("Error starting server: ", err)
		}
	}()

	conn, err := grpc.Dial(s.cfg.ClientsConfig.ProductServiceClientUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		s.logger.Fatalf("Failed to connect: %v\n", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			s.logger.Fatal(err)
		}
	}(conn)

	productServiceClient := pb.NewProductServiceClient(conn)

	//ReadBlaBla(productServiceClient)

	productHandler := handlers.NewProductHandler(s.cfg, s.logger, productServiceClient)

	middlewareManager := middlewares.NewMiddlewareManager(s.cfg, s.logger)
	s.echo.Use(middlewareManager.RequestLoggerMiddleware)

	s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXRequestID},
	}))
	s.echo.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10, //1kb
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	s.echo.Use(middleware.RequestID())
	//s.echo.Use(middlewareManager.MetricsMiddleware(metrics))
	s.echo.Use(middleware.Secure())
	s.echo.Use(middleware.BodyLimit("2M"))
	s.echo.GET("/swagger/*", echoSwagger.WrapHandler)

	v1 := s.echo.Group("/api/v1")
	health := v1.Group("/health")
	productGroup := v1.Group("/products")

	handlers.MapProductRoutes(productGroup, productHandler)

	health.GET("", func(c echo.Context) error {
		s.logger.Infof("Health check RequestID: %s", util.GetRequestId(c))
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	// gracefull shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	ctx, shutdown := context.WithTimeout(context.Background(), time.Duration(s.cfg.Server.CtxTimeout)*time.Second)
	defer shutdown()
	s.logger.Info("Server exited properly")
	return s.echo.Server.Shutdown(ctx)
}

func ReadBlaBla(c pb.ProductServiceClient) {
	log.Println("It's working!")
	res, err := c.GetProductDetail(context.Background(), &pb.GetProductDetailRequest{
		Id: "123",
	})
	if err != nil {
		log.Printf("Error happened!: %v\n", err)
	}

	log.Println(res)
}

func NewServer(cfg *config.Config, logger logger.Logger) *Server {
	return &Server{
		echo:   echo.New(),
		cfg:    cfg,
		logger: logger,
	}
}
