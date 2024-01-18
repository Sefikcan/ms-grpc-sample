package mongo

import (
	"context"
	"fmt"
	"github.com/sefikcan/ms-grpc-sample/product/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func NewMongo(c *config.Config) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*c.Server.CtxTimeout))
	defer cancel()

	connStr := fmt.Sprintf("mongodb://%s:%s", c.Mongo.Host, c.Mongo.Port)

	clientOpts := options.Client().ApplyURI(connStr)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
