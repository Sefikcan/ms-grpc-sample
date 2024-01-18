package use_case

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/sefikcan/ms-grpc-sample/product/internal/entity"
	"github.com/sefikcan/ms-grpc-sample/product/internal/repository"
	"github.com/sefikcan/ms-grpc-sample/product/internal/utils"
	"github.com/sefikcan/ms-grpc-sample/product/pkg/config"
	"github.com/sefikcan/ms-grpc-sample/product/pkg/logger"
	pb "github.com/sefikcan/ms-grpc-sample/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductUseCase interface {
	Create(ctx context.Context, request *pb.CreateProductRequest) (*pb.CreateProductResponse, error)
	GetById(ctx context.Context, request *pb.GetProductDetailRequest) (*pb.GetProductDetailResponse, error)
	Delete(ctx context.Context, request *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error)
	Update(ctx context.Context, request *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error)
}

type productUseCase struct {
	cfg               *config.Config
	productRepository repository.ProductRepository
	logger            logger.Logger
}

func (p productUseCase) Create(ctx context.Context, request *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	span, spanContext := opentracing.StartSpanFromContext(ctx, "productUseCase.Create")
	defer span.Finish()

	product := entity.Product{
		Name:       request.Name,
		Category:   request.Category,
		OptionName: request.OptionName,
	}

	res, err := p.productRepository.Create(spanContext, product)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal Error: %v\n", err),
		)
	}

	return &pb.CreateProductResponse{
		Name:       res.Name,
		Category:   res.Category,
		OptionName: res.OptionName,
		Id:         res.Id.String(),
	}, nil
}

func (p productUseCase) GetById(ctx context.Context, request *pb.GetProductDetailRequest) (*pb.GetProductDetailResponse, error) {
	span, spanContext := opentracing.StartSpanFromContext(ctx, "productUseCase.GetById")
	defer span.Finish()

	oid, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Cannot parse Id",
		)
	}

	res, err := p.productRepository.GetById(spanContext, oid)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal Error: %v\n", err),
		)
	}

	return utils.DocumentToProduct(res), nil
}

func (p productUseCase) Delete(ctx context.Context, request *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	span, spanContext := opentracing.StartSpanFromContext(ctx, "productUseCase.Delete")
	defer span.Finish()

	oid, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Cannot parse Id",
		)
	}

	err = p.productRepository.Delete(spanContext, oid)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal Error: %v\n", err),
		)
	}

	return &pb.DeleteProductResponse{}, nil
}

func (p productUseCase) Update(ctx context.Context, request *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	span, spanContext := opentracing.StartSpanFromContext(ctx, "productUseCase.Update")
	defer span.Finish()

	oid, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Cannot parse Id",
		)
	}

	product := entity.Product{
		Id:         oid,
		Name:       request.Name,
		OptionName: request.OptionName,
		Category:   request.Category,
	}

	res, err := p.productRepository.Update(spanContext, product)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal Error: %v\n", err),
		)
	}

	return &pb.UpdateProductResponse{
		Name:       res.Name,
		Category:   res.Category,
		OptionName: res.OptionName,
		Id:         res.Id.String(),
	}, nil
}

func NewProductUseCase(cfg *config.Config, productRepository repository.ProductRepository, logger logger.Logger) ProductUseCase {
	return &productUseCase{
		cfg:               cfg,
		productRepository: productRepository,
		logger:            logger,
	}
}
