package handler

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gosrv/redis"
	"io"
	"net/http"
	"net/url"
	"time"
)

const base string = "http://localhost:8080"

var basectx context.Context = context.Background()
var client = redis.NewRedisInstance()

type res struct {
	ShortURL string `json:"ShortURL"`
}

func validationURL(input string) bool {
	parsedURI, err := url.ParseRequestURI(input)
	return err == nil && parsedURI.Scheme != "" && parsedURI.Host != ""
}

func generateShortURL(input string) string {
	hash := sha256.Sum256([]byte(input + fmt.Sprint(time.Now().UnixNano())))
	return base64.URLEncoding.EncodeToString(hash[:])[:8]
}

type req struct {
	InputURL string `json:"inURl"`
}

func HandleShortening(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	fmt.Println("here")
	// inputURL := r.PathValue("inputURL")
	bindata, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	js := new(req)
	json.Unmarshal(bindata, &js)
	inputURL := js.InputURL
	if !validationURL(inputURL) {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	var key string = generateShortURL(inputURL)
	sURL := fmt.Sprintf("%s/%s", base, key)

	error := client.Set(basectx, key, inputURL, 0).Err()
	if error != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := res{ShortURL: sURL}
	resBytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resBytes)
}
