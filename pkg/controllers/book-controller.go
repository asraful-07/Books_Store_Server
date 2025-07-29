package controllers

import (
	"bookstore/pkg/config"
	"bookstore/pkg/models"
	"bookstore/pkg/utils"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


var bookCollection = config.GetCollection("books")
var userCollection = config.GetCollection("users")


func RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.User
	utils.ParseBody(r, &user)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if email already exists
	count, _ := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
	if count > 0 {
		http.Error(w, "Email already exists", http.StatusBadRequest)
		return
	}

	// Default role = "user" if not provided
	if user.Role == "" {
		user.Role = "user"
	}

	_, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		http.Error(w, "Failed to register", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Registration successful",
	})
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var loginData models.User
	utils.ParseBody(r, &loginData)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{
		"email":    loginData.Email,
		"password": loginData.Password,
	}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		} else {
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
		return
	}

	// Success – return user info (except password)
	user.Password = ""
	json.NewEncoder(w).Encode(user)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book models.Book
	utils.ParseBody(r, &book)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := bookCollection.InsertOne(ctx, book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	book.ID = res.InsertedID.(primitive.ObjectID)
	json.NewEncoder(w).Encode(book)
}

func GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var books []models.Book

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := bookCollection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var book models.Book
		cursor.Decode(&book)
		books = append(books, book)
	}

	json.NewEncoder(w).Encode(books)
}

func GetRelatedBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	bookID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Step 1: Find the current book by ID
	var currentBook models.Book
	err = bookCollection.FindOne(ctx, bson.M{"_id": bookID}).Decode(&currentBook)
	if err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	// Step 2: Define filter to find related books by productType
	filter := bson.M{
		"productType": currentBook.ProductType,
		"_id":         bson.M{"$ne": bookID}, 
	}

	// Step 3: Limit results to 4 using options
	findOptions := options.Find().SetLimit(12)

	cursor, err := bookCollection.Find(ctx, filter, findOptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var relatedBooks []models.Book
	for cursor.Next(ctx) {
		var book models.Book
		cursor.Decode(&book)
		relatedBooks = append(relatedBooks, book)
	}

	json.NewEncoder(w).Encode(relatedBooks)
}

func GetBookById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["bookID"]
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var book models.Book
	err = bookCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&book)
	if err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(book)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["bookID"]
	objID, _ := primitive.ObjectIDFromHex(id)

	var book models.Book
	utils.ParseBody(r, &book)

	update := bson.M{
		"$set": book,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := bookCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		http.Error(w, "Update failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("Book updated")
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["bookID"]
	objID, _ := primitive.ObjectIDFromHex(id)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := bookCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		http.Error(w, "Delete failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("Book deleted")
}

func GetByBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var books []models.Book

	category := r.URL.Query().Get("category")

	filter := bson.M{}
	if category != "" && category != "All" {
		filter = bson.M{"productType": category}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := bookCollection.Find(ctx, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var book models.Book
		if err := cursor.Decode(&book); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}

	if err := cursor.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(books)
}

