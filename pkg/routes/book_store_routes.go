package routes

import (
	"bookstore/pkg/controllers"

	"github.com/gorilla/mux"
)

func RegisterBookStorRoutes(router *mux.Router) {
	router.HandleFunc("/book", controllers.CreateBook).Methods("POST")
	router.HandleFunc("/books", controllers.GetBook).Methods("GET")
	router.HandleFunc("/book", controllers.GetBooks).Methods("GET")
	router.HandleFunc("/book/related/{id}", controllers.GetRelatedBooks).Methods("GET")
	router.HandleFunc("/book/{bookID}", controllers.GetBookById).Methods("GET")
	router.HandleFunc("/book/{bookID}", controllers.UpdateBook).Methods("PUT")
	router.HandleFunc("/book/{bookID}", controllers.DeleteBook).Methods("DELETE")
	router.HandleFunc("/register", controllers.RegisterUser).Methods("POST")
	router.HandleFunc("/login", controllers.LoginUser).Methods("POST")
	router.HandleFunc("/book-info", controllers.GetByBook).Methods("GET")
	router.HandleFunc("/book-data", controllers.GetByBookData).Methods("GET")
	router.HandleFunc("/book-cart", controllers.CreateAddCart).Methods("POST")
	router.HandleFunc("/carts", controllers.GetCartById).Methods("GET")
	router.HandleFunc("/book-cart", controllers.GetCartByEmail).Methods("GET")
	router.HandleFunc("/cart/{cartID}", controllers.DeleteCart).Methods("DELETE")
	router.HandleFunc("/favorites", controllers.CreateFavorites).Methods("POST")
	router.HandleFunc("/favorites", controllers.GetFavorites).Methods("GET")
	router.HandleFunc("/book-favorites", controllers.GetFavByEmail).Methods("GET")
	router.HandleFunc("/favorites", controllers.DeleteFavorite).Methods("DELETE")

}
