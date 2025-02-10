package main

import (
	"context"
	"fmt"
	"gosrv/handler"
	"gosrv/redis"
	"io"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	client := redis.NewRedisInstance()
	ping, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	} else {
		log.Println("Redis connected:", ping)
	}

	mux.HandleFunc("POST /api/shorten", handler.HandleShortening)
	mux.HandleFunc("GET /{key}", handler.HandleRedirect)
	mux.HandleFunc("GET /test",
		func(writer http.ResponseWriter, request *http.Request) {
			io.WriteString(writer, "Here is a response.")
		})

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
