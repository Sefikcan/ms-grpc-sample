package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	Id         primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `bson:"name,omitempty"`
	Category   string             `bson:"category"`
	OptionName string             `bson:"optionName,omitempty"`
}
