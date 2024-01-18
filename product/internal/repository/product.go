package repository

import (
	"context"
	"github.com/sefikcan/ms-grpc-sample/product/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepository interface {
	Create(ctx context.Context, product entity.Product) (entity.Product, error)
	Update(ctx context.Context, product entity.Product) (entity.Product, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
	GetById(ctx context.Context, id primitive.ObjectID) (entity.Product, error)
}

type productRepository struct {
	db *mongo.Client
}

func (p productRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	collection := p.db.Database("productDb").Collection("product")

	_, err := collection.DeleteOne(ctx, id)
	return err
}

func (p productRepository) Create(ctx context.Context, product entity.Product) (entity.Product, error) {
	collection := p.db.Database("productDb").Collection("product")

	_, err := collection.InsertOne(ctx, product)
	return product, err
}

func (p productRepository) Update(ctx context.Context, product entity.Product) (entity.Product, error) {
	collection := p.db.Database("productDb").Collection("product")

	filter := bson.M{
		"_id": product.Id,
	}
	update := bson.M{
		"$set": bson.M{
			"name":       product.Name,
			"category":   product.Category,
			"optionName": product.OptionName}}

	_, err := collection.UpdateOne(ctx, filter, update)
	return product, err
}

func (p productRepository) GetById(ctx context.Context, id primitive.ObjectID) (entity.Product, error) {
	collection := p.db.Database("productDb").Collection("product")

	var product entity.Product
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	if err != nil {
		return entity.Product{}, err
	}

	return product, nil
}

func NewProductRepository(db *mongo.Client) ProductRepository {
	return &productRepository{
		db: db,
	}
}
