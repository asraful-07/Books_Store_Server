package main

import (
	"bookstore/pkg/config"
	"bookstore/pkg/routes"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	router := mux.NewRouter()
	config.Connect()
	routes.RegisterBookStorRoutes(router)

	// Enable CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, 
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
