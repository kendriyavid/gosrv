package handler

import (
	"net/http"
)

func HandleRedirect(w http.ResponseWriter, r *http.Request) {
	var key string = r.PathValue("key")
	val, err := client.Get(basectx, key).Result()
	if err != nil {
		http.Error(w, "url not present", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, val, http.StatusPermanentRedirect)
}
