package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go_blog_api/internals/database"
	"go_blog_api/internals/handlers"
	"go_blog_api/internals/middleware"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Get configuration from environmnet
	mongoURI := getEnv("MONGODB_URI", "mongodb://localhost:27017")
	dbName := getEnv("DATABASE_NAME", "blog_api")
	collectionName := getEnv("COLLECTION_NAME", "articles")
	port := getEnv("PORT", 8080)
	host := getEnv("HOST", "localhost")

	// connect to database
	if err := database.ConnectDB(mongoURI); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	setupGracefulShutdown()

	// Create article handler
	articleHandler := handlers.NewArticleHandler(dbName, collectionName)

	// Set up router
	router := mux.NewRouter()

	// Apply middleware
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.JSONMiddleware)
	router.Use(middleware.CORSMiddleware)

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"status": "OK",
			"timestamp": time.Now().Format(time.RFC3339),
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	})

	// Article routes

	// Start server
	address := host + ":" + port
	log.Printf("✈️ Server starting on http://%s", address)
	
	if err := http.ListenAndServe(address, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}	

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return  value
	}
	return  fallback
}

// handles os signal for graceful shutdowns
func setupGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func ()  {
		<-c
		log.Println("\n🛑 Shutting down server...")
		database.DisconnectDB()
		os.Exit(0)
	}()
}
