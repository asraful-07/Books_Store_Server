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
		AllowedOrigins:   []string{"http://localhost:5173", "https://auth-test-project-1a397.web.app"}, 
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)
    port := `:PORT`
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(port, handler))
}
