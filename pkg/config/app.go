package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func Connect() {
	clientOptions := options.Client().ApplyURI(`mongodb+srv://DB_USERNAME:DB_PASSWORD@cluster0.w0mh3.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0`)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("MongoDB ping failed: %v", err)
	}

	// Use your database name here (replace `bookstore`)
	DB = client.Database(`DB_NAME`)
	fmt.Println("Connected to MongoDB!")
}

func GetCollection(collectionName string) *mongo.Collection {
	if DB == nil {
		log.Fatal("DB is nil, make sure Connect() was called.")
	}
	return DB.Collection(collectionName)
}
