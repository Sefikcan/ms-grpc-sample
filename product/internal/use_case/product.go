package use_case

import (
	"context"
	"fmt"
	"github.com/sefikcan/ms-grpc-sample/product/internal/entity"
	"github.com/sefikcan/ms-grpc-sample/product/internal/mappers"
	"github.com/sefikcan/ms-grpc-sample/product/internal/repository"
	"github.com/sefikcan/ms-grpc-sample/product/pkg/config"
	"github.com/sefikcan/ms-grpc-sample/product/pkg/logger"
	pb "github.com/sefikcan/ms-grpc-sample/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type ProductUseCase struct {
	cfg               *config.Config
	productRepository repository.ProductRepository
	logger            logger.Logger
	pb.UnimplementedProductServiceServer
}

func (p ProductUseCase) Create(ctx context.Context, request *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	product := entity.Product{
		Name:       request.Name,
		Category:   request.Category,
		OptionName: request.OptionName,
	}

	res, err := p.productRepository.Create(ctx, product)
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
		Id:         res.Id.Hex(),
	}, nil
}

func (p ProductUseCase) GetById(ctx context.Context, request *pb.GetProductDetailRequest) (*pb.GetProductDetailResponse, error) {
	oid, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Cannot parse Id",
		)
	}

	res, err := p.productRepository.GetById(ctx, oid)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal Error: %v\n", err),
		)
	}

	return mappers.DocumentToProduct(res), nil
}

func (p ProductUseCase) Delete(ctx context.Context, request *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	oid, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Cannot parse Id",
		)
	}

	err = p.productRepository.Delete(ctx, oid)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal Error: %v\n", err),
		)
	}

	return &pb.DeleteProductResponse{}, nil
}

func (p ProductUseCase) Update(ctx context.Context, request *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
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

	res, err := p.productRepository.Update(ctx, product)
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
		Id:         res.Id.Hex(),
	}, nil
}

type ProductServerStruct struct {
	pb.UnimplementedProductServiceServer
	productUseCase *ProductUseCase
}

func NewProductUseCase(cfg *config.Config, productRepository repository.ProductRepository, logger logger.Logger, grpcServer *grpc.Server) *ProductUseCase {
	productGrpc := &ProductServerStruct{
		productUseCase: &ProductUseCase{
			cfg:               cfg,
			productRepository: productRepository,
			logger:            logger,
		},
	}
	pb.RegisterProductServiceServer(grpcServer, productGrpc)
	return productGrpc.productUseCase
}

func (s *ProductServerStruct) CreateProduct(ctx context.Context, in *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	createdProduct, err := s.productUseCase.Create(ctx, in)
	if err != nil {
		log.Printf("Failed to create product: %v\n", err)
		return nil, err
	}

	return createdProduct, nil
}

func (s *ProductServerStruct) UpdateProduct(ctx context.Context, in *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	updatedProduct, err := s.productUseCase.Update(ctx, in)
	if err != nil {
		log.Printf("Failed to update product: %v\n", err)
		return nil, err
	}

	return updatedProduct, nil
}

func (s *ProductServerStruct) DeleteProduct(ctx context.Context, in *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	deletedProduct, err := s.productUseCase.Delete(ctx, in)
	if err != nil {
		log.Printf("Failed to delete product: %v\n", err)
		return nil, err
	}

	return deletedProduct, nil
}

func (s *ProductServerStruct) GetProductDetail(ctx context.Context, in *pb.GetProductDetailRequest) (*pb.GetProductDetailResponse, error) {
	product, err := s.productUseCase.GetById(ctx, in)
	if err != nil {
		log.Printf("Failed to get product: %v\n", err)
		return nil, err
	}

	return product, nil
}
