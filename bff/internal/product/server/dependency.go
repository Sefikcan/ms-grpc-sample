package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sefikcan/ms-grpc-sample/bff/internal/product/handlers"
	"github.com/sefikcan/ms-grpc-sample/bff/internal/product/middlewares"
	"github.com/sefikcan/ms-grpc-sample/bff/pkg/metric"
	"github.com/sefikcan/ms-grpc-sample/bff/pkg/util"
	pb "github.com/sefikcan/ms-grpc-sample/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
)

func (s *Server) MapHandlers(e *echo.Echo) error {
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
