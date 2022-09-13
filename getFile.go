package main

import (
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getFile(w http.ResponseWriter, r *http.Request) {
	_id := strings.TrimPrefix(r.URL.Path, "/assets/")
	objectId, err := primitive.ObjectIDFromHex(_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := Get(objectId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Add("Cache-Control", "max-age=2592000")
	http.ServeFile(w, r, result.FilePath)
}
