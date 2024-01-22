package client

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/sefikcan/ms-grpc-sample/bff/pkg/logger"
	pb "github.com/sefikcan/ms-grpc-sample/proto"
)

type ProductClient interface {
	GetById(ctx context.Context, id string, client pb.ProductServiceClient) (*pb.GetProductDetailResponse, error)
	Create(ctx context.Context, request *pb.CreateProductRequest, client pb.ProductServiceClient) (*pb.CreateProductResponse, error)
	Delete(ctx context.Context, request *pb.DeleteProductRequest, client pb.ProductServiceClient) (*pb.DeleteProductResponse, error)
	Update(ctx context.Context, request *pb.UpdateProductRequest, client pb.ProductServiceClient) (*pb.UpdateProductResponse, error)
}

type productClient struct {
	logger logger.Logger
}

func (p productClient) Update(ctx context.Context, request *pb.UpdateProductRequest, client pb.ProductServiceClient) (*pb.UpdateProductResponse, error) {
	span, spanContext := opentracing.StartSpanFromContext(ctx, "productClient.Update")
	defer span.Finish()

	res, err := client.UpdateProduct(spanContext, request)
	if err != nil {
		p.logger.Errorf("Product Update Error: %v\n", err)
		return nil, err
	}

	p.logger.Infof("Product was updated. Id: %s\n", res.Id)

	return res, nil
}

func (p productClient) Delete(ctx context.Context, request *pb.DeleteProductRequest, client pb.ProductServiceClient) (*pb.DeleteProductResponse, error) {
	span, spanContext := opentracing.StartSpanFromContext(ctx, "productClient.Delete")
	defer span.Finish()

	_, err := client.DeleteProduct(spanContext, request)
	if err != nil {
		p.logger.Errorf("Product Delete Error: %v\n", err)
		return nil, err
	}

	p.logger.Info("Product was deleted!")

	return &pb.DeleteProductResponse{}, nil
}

func (p productClient) GetById(ctx context.Context, id string, client pb.ProductServiceClient) (*pb.GetProductDetailResponse, error) {
	span, spanContext := opentracing.StartSpanFromContext(ctx, "productClient.GetById")
	defer span.Finish()

	res, err := client.GetProductDetail(spanContext, &pb.GetProductDetailRequest{
		Id: id,
	})
	if err != nil {
		p.logger.Errorf("Product Get Detail Error: %v\n", err)
		return nil, err
	}

	p.logger.Infof("Product was read: %v\n", res)

	return res, nil
}

func (p productClient) Create(ctx context.Context, request *pb.CreateProductRequest, client pb.ProductServiceClient) (*pb.CreateProductResponse, error) {
	span, spanContext := opentracing.StartSpanFromContext(ctx, "productClient.Create")
	defer span.Finish()

	res, err := client.CreateProduct(spanContext, request)
	if err != nil {
		p.logger.Errorf("Product Create Error: %v\n", err)
		return nil, err
	}

	p.logger.Infof("Product has been created: %s\n", res.Id)

	return res, nil
}

func NewProductClient(logger logger.Logger) ProductClient {
	return &productClient{
		logger: logger,
	}
}
