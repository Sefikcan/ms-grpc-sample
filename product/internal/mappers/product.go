package mappers

import (
	"github.com/sefikcan/ms-grpc-sample/product/internal/entity"
	pb "github.com/sefikcan/ms-grpc-sample/proto"
)

func DocumentToProduct(data entity.Product) *pb.GetProductDetailResponse {
	return &pb.GetProductDetailResponse{
		Id:         data.Id.Hex(),
		Name:       data.Name,
		OptionName: data.OptionName,
		Category:   data.Category,
	}
}
