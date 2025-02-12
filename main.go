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
	basectx := context.Background()
	shortener := handler.NewURLshortener(client, "http://localhost:8080")
	decompressor := handler.NewURLDecompressor(client)
	ping, err := client.Ping(basectx).Result()
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	} else {
		log.Println("Redis connected:", ping)
	}

	mux.HandleFunc("POST /api/shorten", shortener.HandleShortening)
	mux.HandleFunc("GET /{key}", decompressor.HandleRedirect)
	mux.HandleFunc("GET /test",
		func(writer http.ResponseWriter, request *http.Request) {
			io.WriteString(writer, "Here is a response.")
		})

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
