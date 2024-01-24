package mappers

import (
	"github.com/sefikcan/ms-grpc-sample/bff/internal/product/dto/requests"
	"github.com/sefikcan/ms-grpc-sample/bff/internal/product/dto/responses"
	pb "github.com/sefikcan/ms-grpc-sample/proto"
)

func CreateProductRequestToGrpcRequestObject(productRequest requests.CreateProductRequest) *pb.CreateProductRequest {
	return &pb.CreateProductRequest{
		Name:       productRequest.Name,
		Category:   productRequest.Category,
		OptionName: productRequest.OptionName,
	}
}

func UpdateProductRequestToGrpcRequestObject(id string, productRequest requests.UpdateProductRequest) *pb.UpdateProductRequest {
	return &pb.UpdateProductRequest{
		Id:         id,
		Name:       productRequest.Name,
		Category:   productRequest.Category,
		OptionName: productRequest.OptionName,
	}
}

func CreateProductGrpcResponseToResponseObject(productResponse *pb.CreateProductResponse) *responses.ProductResponse {
	return &responses.ProductResponse{
		Id:         productResponse.Id,
		Name:       productResponse.Name,
		OptionName: productResponse.OptionName,
		Category:   productResponse.Category,
	}
}

func GetProductGrpcResponseToResponseObject(productResponse *pb.GetProductDetailResponse) *responses.ProductResponse {
	return &responses.ProductResponse{
		Id:         productResponse.Id,
		Name:       productResponse.Name,
		OptionName: productResponse.OptionName,
		Category:   productResponse.Category,
	}
}
