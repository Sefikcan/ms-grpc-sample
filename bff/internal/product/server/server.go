package server

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sefikcan/ms-grpc-sample/bff/internal/product/handlers"
	"github.com/sefikcan/ms-grpc-sample/bff/internal/product/middlewares"
	"github.com/sefikcan/ms-grpc-sample/bff/pkg/config"
	"github.com/sefikcan/ms-grpc-sample/bff/pkg/logger"
	"github.com/sefikcan/ms-grpc-sample/bff/pkg/metric"
	"github.com/sefikcan/ms-grpc-sample/bff/pkg/util"
	pb "github.com/sefikcan/ms-grpc-sample/proto"
	echoSwagger "github.com/swaggo/echo-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	if err := s.Register(s.echo); err != nil {
		return err
	}

	// gracefull shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	ctx, shutdown := context.WithTimeout(context.Background(), time.Duration(s.cfg.Server.CtxTimeout)*time.Second)
	defer shutdown()
	s.logger.Info("Server exited properly")
	return s.echo.Server.Shutdown(ctx)
}

func (s *Server) Register(e *echo.Echo) error {
	metrics, err := metric.CreateMetrics(s.cfg.Metric.Url, s.cfg.Metric.ServiceName)
	if err != nil {
		s.logger.Errorf("CreateMetrics error: %s", err)
	}
	s.logger.Infof("Metrics available URL: %s, ServiceName: %s", s.cfg.Metric.Url, s.cfg.Metric.ServiceName)

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

	productHandler := handlers.NewProductHandler(s.cfg, s.logger, productServiceClient)

	middlewareManager := middlewares.NewMiddlewareManager(s.cfg, s.logger)
	e.Use(middlewareManager.RequestLoggerMiddleware)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXRequestID},
	}))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10, //1kb
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	e.Use(middleware.RequestID())
	e.Use(middlewareManager.MetricsMiddleware(metrics))
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("2M"))
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	v1 := e.Group("/api/v1")
	health := v1.Group("/health")
	productGroup := v1.Group("/products")

	handlers.MapProductRoutes(productGroup, productHandler)

	health.GET("", func(c echo.Context) error {
		s.logger.Infof("Health check RequestID: %s", util.GetRequestId(c))
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	return nil
}

func NewServer(cfg *config.Config, logger logger.Logger) *Server {
	return &Server{
		echo:   echo.New(),
		cfg:    cfg,
		logger: logger,
	}
}
