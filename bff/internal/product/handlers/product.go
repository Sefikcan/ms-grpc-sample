package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/sefikcan/ms-grpc-sample/bff/internal/product/dto/requests"
	"github.com/sefikcan/ms-grpc-sample/bff/internal/product/mappers"
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

// Create godoc
// @Summary Create product
// @Description Create product handler
// @Tags Product
// @Accept json
// @Produce json
// @Param createProductRequest body requests.CreateProductRequest true "Create Product"
// @Success 201 {object} responses.ProductResponse
// @Router /products [post]
func (p productHandlers) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, spanContext := opentracing.StartSpanFromContext(util.GetRequestCtx(c), "productHandler.Create")
		defer span.Finish()

		productRequest := requests.CreateProductRequest{}
		if err := c.Bind(&productRequest); err != nil {
			p.logger.Errorf("Error, RequestId: %s, IPAddress: %s, Error: %s", util.GetRequestId(c), util.GetIPAddress(c), err)
			return c.JSON(http.StatusBadRequest, util.NewHttpResponse(http.StatusBadRequest, strings.ToLower(err.Error()), nil))
		}

		clientReq := mappers.CreateProductRequestToGrpcRequestObject(productRequest)

		res, err := p.productServiceClient.CreateProduct(spanContext, clientReq)
		if err != nil {
			p.logger.Fatalf("Unexpected error: %v\n", err)
		}

		return c.JSON(http.StatusCreated, mappers.CreateProductGrpcResponseToResponseObject(res))
	}
}

// Delete godoc
// @Summary Delete product
// @Description Delete by id product handler
// @Tags Product
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 204
// @Router /products/{id} [delete]
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

// Update godoc
// @Summary Update product
// @Description Update product handler
// @Tags Product
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param updateProductRequest body requests.UpdateProductRequest true "Update Product"
// @Success 200 {object} responses.ProductResponse
// @Router /products/{id} [put]
func (p productHandlers) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, spanContext := opentracing.StartSpanFromContext(util.GetRequestCtx(c), "productHandler.Update")
		defer span.Finish()

		req := requests.UpdateProductRequest{}
		if err := c.Bind(&req); err != nil {
			p.logger.Errorf("Error, RequestId: %s, IPAddress: %s, Error: %s", util.GetRequestId(c), util.GetIPAddress(c), err)
			return c.JSON(http.StatusBadRequest, util.NewHttpResponse(http.StatusBadRequest, strings.ToLower(err.Error()), nil))
		}

		res, err := p.productServiceClient.UpdateProduct(spanContext, mappers.UpdateProductRequestToGrpcRequestObject(req))
		if err != nil {
			p.logger.Errorf("Error happened while updating: %v\n", err)
		}

		return c.JSON(http.StatusOK, res)
	}
}

// GetById godoc
// @Summary Get by id product
// @Description Get by id product handler
// @Tags Product
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} responses.ProductResponse
// @Router /products/{id} [get]
func (p productHandlers) GetById() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, spanContext := opentracing.StartSpanFromContext(util.GetRequestCtx(c), "productHandler.GetById")
		defer span.Finish()

		req := &pb.GetProductDetailRequest{Id: c.Param("id")}

		res, err := p.productServiceClient.GetProductDetail(spanContext, req)
		if err != nil {
			p.logger.Errorf("Error happened while reading: %v\n", err)
		}

		return c.JSON(http.StatusOK, mappers.GetProductGrpcResponseToResponseObject(res))
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
