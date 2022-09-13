package main

import (
	"errors"
	"log"
	"net/http"
	"os"
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
		log.Default().Println("db not found")
		http.Error(w, "no data", http.StatusNotFound)
		os.Remove(result.FilePath)
		return
	}
	w.Header().Add("Cache-Control", "max-age=2592000")
	if _, err := os.Stat(result.FilePath); errors.Is(err, os.ErrNotExist) {
		log.Default().Println("file not found")
		http.Error(w, "no data", http.StatusNotFound)
		Del(objectId)
		return
	}
	http.ServeFile(w, r, result.FilePath)
}
