package handler

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gosrv/redis"
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
	js := new(req)
	var decoder *json.Decoder = json.NewDecoder(r.Body)
	if err := decoder.Decode(js); err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	inputURL := js.InputURL
	if !validationURL(inputURL) {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	var key string = generateShortURL(inputURL)
	sURL := fmt.Sprintf("%s/%s", base, key)

	if err := client.Set(basectx, key, inputURL, 0).Err(); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := res{ShortURL: sURL}
	var encoder json.Encoder = *json.NewEncoder(w)

	if err := encoder.Encode(response); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

}
