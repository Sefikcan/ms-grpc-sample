package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/sefikcan/ms-grpc-sample/bff/pkg/config"
	"github.com/sefikcan/ms-grpc-sample/bff/pkg/logger"
	"github.com/sefikcan/ms-grpc-sample/bff/pkg/util"
	pb "github.com/sefikcan/ms-grpc-sample/proto"
	"net/http"
	"strings"
)

type ProductHandlers interface {
	Create() echo.HandlerFunc
	Delete() echo.HandlerFunc
	Update() echo.HandlerFunc
	GetById() echo.HandlerFunc
	GetAll() echo.HandlerFunc
}

type productHandlers struct {
	cfg                  *config.Config
	logger               logger.Logger
	productServiceClient pb.ProductServiceClient
}

func (p productHandlers) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, spanContext := opentracing.StartSpanFromContext(util.GetRequestCtx(c), "productHandler.Create")
		defer span.Finish()

		productRequest := &pb.CreateProductRequest{}
		if err := c.Bind(productRequest); err != nil {
			p.logger.Errorf("Error, RequestId: %s, IPAddress: %s, Error: %s", util.GetRequestId(c), util.GetIPAddress(c), err)
			return c.JSON(http.StatusBadRequest, util.NewHttpResponse(http.StatusBadRequest, strings.ToLower(err.Error()), nil))
		}

		res, err := p.productServiceClient.CreateProduct(spanContext, productRequest)
		if err != nil {
			p.logger.Fatalf("Unexpected error: %v\n", err)
		}

		return c.JSON(http.StatusCreated, res)
	}
}

func (p productHandlers) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, spanContext := opentracing.StartSpanFromContext(util.GetRequestCtx(c), "productHandler.Delete")
		defer span.Finish()

		req := &pb.DeleteProductRequest{
			Id: c.Param("id"),
		}

		_, err := p.productServiceClient.DeleteProduct(spanContext, req)
		if err != nil {
			p.logger.Errorf("Error happened while deleting: %v\n", err)
		}

		return c.NoContent(http.StatusNoContent)
	}
}

func (p productHandlers) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, spanContext := opentracing.StartSpanFromContext(util.GetRequestCtx(c), "productHandler.Update")
		defer span.Finish()

		req := &pb.UpdateProductRequest{}
		if err := c.Bind(req); err != nil {
			p.logger.Errorf("Error, RequestId: %s, IPAddress: %s, Error: %s", util.GetRequestId(c), util.GetIPAddress(c), err)
			return c.JSON(http.StatusBadRequest, util.NewHttpResponse(http.StatusBadRequest, strings.ToLower(err.Error()), nil))
		}

		res, err := p.productServiceClient.UpdateProduct(spanContext, req)
		if err != nil {
			p.logger.Errorf("Error happened while updating: %v\n", err)
		}

		return c.JSON(http.StatusOK, res)
	}
}

func (p productHandlers) GetById() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, spanContext := opentracing.StartSpanFromContext(util.GetRequestCtx(c), "productHandler.GetById")
		defer span.Finish()

		req := &pb.GetProductDetailRequest{Id: c.Param("id")}

		res, err := p.productServiceClient.GetProductDetail(spanContext, req)
		if err != nil {
			p.logger.Errorf("Error happened while reading: %v\n", err)
		}

		return c.JSON(http.StatusOK, res)
	}
}

func (p productHandlers) GetAll() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusNoContent, nil)
	}
}

func NewProductHandler(cfg *config.Config, logger logger.Logger, productServiceClient pb.ProductServiceClient) ProductHandlers {
	return &productHandlers{
		cfg:                  cfg,
		logger:               logger,
		productServiceClient: productServiceClient,
	}
}
