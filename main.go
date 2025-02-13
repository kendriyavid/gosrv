// package main

// import (
// 	"context"
// 	"fmt"
// 	"gosrv/handler"
// 	"gosrv/redis"
// 	"io"
// 	"log"
// 	"net/http"
// )

// func main() {
// 	mux := http.NewServeMux()
// 	client := redis.NewRedisInstance()
// 	basectx := context.Background()
// 	shortener := handler.NewURLshortener(client, "http://localhost:8080")
// 	decompressor := handler.NewURLDecompressor(client)
// 	ping, err := client.Ping(basectx).Result()
// 	if err != nil {
// 		log.Fatalf("Redis connection failed: %v", err)
// 	} else {
// 		log.Println("Redis connected:", ping)
// 	}

// 	mux.HandleFunc("POST /api/shorten", shortener.HandleShortening)
// 	mux.HandleFunc("GET /{key}", decompressor.HandleRedirect)
// 	mux.HandleFunc("GET /test",
// 		func(writer http.ResponseWriter, request *http.Request) {
// 			io.WriteString(writer, "Here is a response.")
// 		})

// 	fmt.Println("Server is running on http://localhost:8080")
// 	log.Fatal(http.ListenAndServe(":8080", mux))
// }

package main

import (
	"context"
	"fmt"
	"gosrv/handler"
	"gosrv/redis"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	_ = godotenv.Load()
	// Load port and base URL from environment variables (with defaults)
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080" // Default port
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:" + port // Default base URL
	}

	// Initialize Redis client
	mux := http.NewServeMux()
	client := redis.NewRedisInstance()
	basectx := context.Background()
	shortener := handler.NewURLshortener(client, baseURL)
	decompressor := handler.NewURLDecompressor(client)

	// Check Redis connection
	ping, err := client.Ping(basectx).Result()
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	} else {
		log.Println("Redis connected:", ping)
	}

	// Define routes
	mux.HandleFunc("POST /api/shorten", shortener.HandleShortening)
	mux.HandleFunc("GET /{key}", decompressor.HandleRedirect)
	mux.HandleFunc("GET /test",
		func(writer http.ResponseWriter, request *http.Request) {
			io.WriteString(writer, "Here is a response.")
		})

	// Start server

	// Enable CORS
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Set this to your frontend origin in production
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	}).Handler(mux)

	serverAddr := ":" + port
	fmt.Printf("Server is running on %s\n", baseURL)
	log.Fatal(http.ListenAndServe(serverAddr, handler))
}
