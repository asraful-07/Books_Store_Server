package models

import (
	"bookstore/pkg/config"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
	Role     string             `json:"role,omitempty" bson:"role,omitempty"` 
}

type Book struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name         string             `json:"name,omitempty" bson:"name,omitempty"`
	Author       string             `json:"author,omitempty" bson:"author,omitempty"`
	Publication  string             `json:"publication,omitempty" bson:"publication,omitempty"`
	ProductType  string             `json:"productType,omitempty" bson:"productType,omitempty"`
	PublishDate  string             `json:"publishDate,omitempty" bson:"publishDate,omitempty"`
	// Color        []string           `json:"color,omitempty" bson:"color,omitempty"` 
	// Size         string             `json:"size,omitempty" bson:"size,omitempty"`  
	// Material     []string           `json:"material,omitempty" bson:"material,omitempty"`
	// Quantity     int                `json:"quantity,omitempty" bson:"quantity,omitempty"`
	Price        float64            `json:"price,omitempty" bson:"price,omitempty"`
	OldPrice     float64            `json:"oldPrice,omitempty" bson:"oldPrice,omitempty"`
	Stock        int                `json:"stock,omitempty" bson:"stock,omitempty"` 
	Description  string             `json:"description,omitempty" bson:"description,omitempty"`
	ImageURLs    []string           `json:"imageUrls,omitempty" bson:"imageUrls,omitempty"` 
	Tags         []string           `json:"tags,omitempty" bson:"tags,omitempty"` 
}

type Cart struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"` 
	Name        string             `json:"name,omitempty" bson:"name,omitempty"` 
	Author      string             `json:"author,omitempty" bson:"author,omitempty"` 
	Price       float64            `json:"price,omitempty" bson:"price,omitempty"` 
	ProductType string             `json:"productType,omitempty" bson:"productType,omitempty"`
	Quantity    int                `json:"quantity,omitempty" bson:"quantity,omitempty"` 
	ImageUrls   []string           `json:"imageUrls,omitempty" bson:"imageUrls,omitempty"` 
	UserEmail   string             `json:"userEmail,omitempty" bson:"userEmail,omitempty"` 
}

func init() {
	config.Connect()
}
