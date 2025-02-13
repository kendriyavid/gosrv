package handler

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	rcli "gosrv/redis"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/klauspost/compress/zstd"
	"github.com/redis/go-redis/v9"
)

// const base string = "http://localhost:3000"
const minLengthCompress int = 100

// _ = godotenv.Load()
var client = rcli.NewRedisInstance()

type URLshortener struct {
	base           string
	client         *redis.Client
	reqPool        sync.Pool
	resPool        sync.Pool
	bufPool        sync.Pool
	compressorPool sync.Pool
	urlTTL         time.Duration
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
		client: redisClient,
		base:   baseURL,
		urlTTL: 24 * time.Hour,
	}
	shortener.reqPool.New = func() interface{} {
		return new(req)
	}

	shortener.resPool.New = func() interface{} {
		return new(res)
	}
	shortener.bufPool.New = func() interface{} {
		return new(bytes.Buffer)
	}
	shortener.compressorPool.New = func() interface{} {
		enc, err := zstd.NewWriter(nil)
		if err != nil {
			log.Fatal(err)
		}
		return enc
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
	fmt.Println(inputURL)
	if !validationURL(inputURL) {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	js.reqOBreset()
	us.reqPool.Put(js)

	var key string = generateShortURL(inputURL)
	sURL := fmt.Sprintf("%s/%s", us.base, key)

	if len(inputURL) >= minLengthCompress {
		// do the compression
		temp := us.bufPool.Get().(*bytes.Buffer)
		temp.Reset()
		enc := us.compressorPool.Get().(*zstd.Encoder)
		enc.Reset(temp)
		if _, err := enc.Write([]byte(inputURL)); err != nil {
			log.Println(err)
			http.Error(w, "encoding problem", http.StatusInternalServerError)
			return
		}
		enc.Close()
		us.compressorPool.Put(enc)

		if err := client.Set(r.Context(), key, append(temp.Bytes(), 1), us.urlTTL).Err(); err != nil {
			fmt.Println("here is the problem")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		us.bufPool.Put(temp)
	} else {
		// appending 0 when not compressed
		if err := client.Set(r.Context(), key, append([]byte(inputURL), byte(0)), us.urlTTL).Err(); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
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
