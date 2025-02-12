package handler

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/klauspost/compress/zstd"
	"github.com/redis/go-redis/v9"
)

type URLdecompressor struct {
	decompressorPool sync.Pool
	client           *redis.Client
}

func NewURLDecompressor(redisClient *redis.Client) *URLdecompressor {
	decompressor := &URLdecompressor{
		client: redisClient,
	}
	decompressor.decompressorPool.New = func() interface{} {
		dec, err := zstd.NewReader(nil)
		if err != nil {
			log.Fatal(err)
		}
		return dec
	}
	return decompressor
}

func (dc *URLdecompressor) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	log.Printf("Processing key: %s", key)

	// Get the stored value as bytes
	val, err := dc.client.Get(r.Context(), key).Bytes()
	if err != nil {
		log.Printf("Failed to get URL from Redis: %v", err)
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}
	log.Printf("Stored value: %v", val[:10]) // Print first 10 bytes

	isCompressed := val[len(val)-1]

	var url string
	if isCompressed == 1 {
		// Get a decoder from the pool
		decoder := dc.decompressorPool.Get().(*zstd.Decoder)
		defer dc.decompressorPool.Put(decoder)

		// Reset the decoder with the new input
		decoder.Reset(bytes.NewReader(val[:len(val)-1]))

		// Read all decompressed data
		decompressed, err := io.ReadAll(decoder)
		if err != nil {
			log.Printf("Decompression failed: %v", err)
			http.Error(w, "Failed to decompress URL", http.StatusInternalServerError)
			return
		}

		url = string(decompressed)
		log.Printf("Decompressed URL: %s", url)
	} else {
		url = string(val[:len(val)-1]) // Convert bytes to string for uncompressed URLs, excluding the flag
	}

	if !validationURL(url) {
		log.Printf("URL validation failed for: %s", url)
		http.Error(w, "Invalid stored URL", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url, http.StatusPermanentRedirect)
}
