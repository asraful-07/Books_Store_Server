package controllers

import (
	"bookstore/pkg/config"
	"bookstore/pkg/models"
	"bookstore/pkg/utils"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


var bookCollection = config.GetCollection("books")
var userCollection = config.GetCollection("users")
var cartsCollection = config.GetCollection("carts")
var favoritesCollection = config.GetCollection("favorites")



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

func GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Query params
	category := r.URL.Query().Get("category")
	search := r.URL.Query().Get("search")
	sortBy := r.URL.Query().Get("sort") 
	bookType := r.URL.Query().Get("type") 
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	filter := bson.M{}

	if category != "" {
		filter["productType"] = category
	}
	if search != "" {
		filter["name"] = bson.M{"$regex": search, "$options": "i"}
	}
	if bookType != "" {
		if bookType == "new" {
			filter["publishDate"] = bson.M{"$gte": "2023"} 
		} else {
			filter["publishDate"] = bson.M{"$lt": "2023"}
		}
	}

	// Pagination
	page, _ := strconv.Atoi(pageStr)
	if page == 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(limitStr)
	if limit == 0 {
		limit = 9
	}
	skip := (page - 1) * limit

	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))

	// Sorting
	switch sortBy {
    case "low":
		findOptions.SetSort(bson.D{{Key: "price", Value: 1}})
	case "high":
		findOptions.SetSort(bson.D{{Key: "price", Value: -1}})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := bookCollection.Find(ctx, filter, findOptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var books []models.Book
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

// This is Real CRUD Operation ***************

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

func GetBooks(w http.ResponseWriter, r *http.Request) {
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

// This is Real CRUD Operation ***************

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

func GetByBookData(w http.ResponseWriter, r *http.Request) {
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

// Cart to Book

func CreateAddCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var cart models.Cart
	utils.ParseBody(r, &cart)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := cartsCollection.InsertOne(ctx, cart)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cart.ID = res.InsertedID.(primitive.ObjectID)
	json.NewEncoder(w).Encode(cart)

}

func GetCartById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") 
	var carts []models.Cart

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := cartsCollection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var cart models.Cart
		cursor.Decode(&cart)
		carts = append(carts, cart)
	}
	json.NewEncoder(w).Encode(carts)
}

func GetCartByEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// email from query param: /book-cart?email=someone@gmail.com
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email query param is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"userEmail": email}

	cursor, err := cartsCollection.Find(ctx, filter)
	if err != nil {
		http.Error(w, "Error fetching cart: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var userCarts []models.Cart
	if err := cursor.All(ctx, &userCarts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(userCarts)
}

func DeleteCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") 
	id := mux.Vars(r)["cartID"]
	objId, _ := primitive.ObjectIDFromHex(id)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := cartsCollection.DeleteOne(ctx, bson.M{"_id" : objId})
	if err != nil {
		http.Error(w, "Delete failed", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode("Book deleted")
}

//  Favorites to Book 

func CreateFavorites(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") 

	var favorite models.Cart
	utils.ParseBody(r, &favorite)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"name":      favorite.Name,
		"userEmail": favorite.UserEmail,
	}
	var existing models.Cart
	err := favoritesCollection.FindOne(ctx, filter).Decode(&existing)
	if err == nil {
	
		http.Error(w, "Book already in wishlist", http.StatusConflict)
		return
	}

	res, err := favoritesCollection.InsertOne(ctx, favorite)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	favorite.ID = res.InsertedID.(primitive.ObjectID)
	json.NewEncoder(w).Encode(favorite)
}

func GetFavorites(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") 
	var favorites []models.Cart

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := favoritesCollection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var favorite models.Cart
		cursor.Decode(&favorite)
		favorites = append(favorites, favorite)
	}
	json.NewEncoder(w).Encode(favorites)
}

func GetFavByEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// email from query param: /book-favorites?email=someone@gmail.com
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email query param is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"userEmail": email}

	cursor, err := favoritesCollection.Find(ctx, filter)
	if err != nil {
		http.Error(w, "Error fetching cart: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var userCarts []models.Cart
	if err := cursor.All(ctx, &userCarts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(userCarts)
}

func DeleteFavorite(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := favoritesCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if res.DeletedCount == 0 {
		http.Error(w, "Favorite not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Deleted successfully"})
}
