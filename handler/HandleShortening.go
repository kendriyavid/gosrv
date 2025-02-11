package handler

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	rcli "gosrv/redis"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

const base string = "http://localhost:8080"

var basectx context.Context = context.Background()
var client = rcli.NewRedisInstance()

type URLshortener struct {
	base    string
	basectx context.Context
	client  *redis.Client
	reqPool sync.Pool
	resPool sync.Pool
	urlTTL  time.Duration
}

type res struct {
	ShortURL string `json:"ShortURL"`
}

type req struct {
	InputURL string `json:"inURl"`
}

func (reqOb *req) reqOBreset() {
	reqOb.InputURL = ""
}

func (resOb *res) resObreset() {
	resOb.ShortURL = ""
}

func NewURLshortener(redisClient *redis.Client, baseURL string) *URLshortener {
	shortener := &URLshortener{
		client:  redisClient,
		base:    baseURL,
		urlTTL:  24 * time.Hour,
		basectx: context.Background(),
	}
	shortener.reqPool.New = func() interface{} {
		return new(req)
	}

	shortener.resPool.New = func() interface{} {
		return new(res)
	}
	return shortener
}

func validationURL(input string) bool {
	parsedURI, err := url.ParseRequestURI(input)
	return err == nil && parsedURI.Scheme != "" && parsedURI.Host != ""
}

func generateShortURL(input string) string {
	hash := sha256.Sum256([]byte(input + fmt.Sprint(time.Now().UnixNano())))
	return base64.URLEncoding.EncodeToString(hash[:])[:8]
}

func (us *URLshortener) HandleShortening(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	js := us.reqPool.Get().(*req)
	if err := json.NewDecoder(r.Body).Decode(js); err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	inputURL := js.InputURL
	if !validationURL(inputURL) {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	js.reqOBreset()
	us.reqPool.Put(js)

	var key string = generateShortURL(inputURL)
	sURL := fmt.Sprintf("%s/%s", base, key)

	if err := client.Set(basectx, key, inputURL, 0).Err(); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := us.resPool.Get().(*res)
	response.ShortURL = sURL

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	response.resObreset()
	us.resPool.Put(response)
}
