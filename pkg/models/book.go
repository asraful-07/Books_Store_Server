package models

import (
	"bookstore/pkg/config"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Author      string             `json:"author,omitempty" bson:"author,omitempty"`
	Publication string             `json:"publication,omitempty" bson:"publication,omitempty"`
}


func init() {
	config.Connect()
}
