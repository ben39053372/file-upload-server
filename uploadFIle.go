package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FilesData struct {
	Size      int64     `bson:"size,omitempty"`
	CreatedAt time.Time `bson:"createdAt,omitempty"`
	FilePath  string    `bson:"filePath,omitempty"`
	Url       string    `bson:"url,omitempty"`
	Type      string    `bson:"type,omitempty"`
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	// Maximum upload of 10 MB files
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	assetsPath := "public/assets/" + time.Now().Format("2006-01-02")

	// create folder
	if _, err := os.Stat(assetsPath); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist
		if err := os.MkdirAll(assetsPath, os.ModePerm); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// create empty doc for get id
	id, err := CreateEmptyDoc()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// gen file path
	fp := assetsPath + "/" + id.(primitive.ObjectID).Hex() + filepath.Ext(handler.Filename)

	// Create file
	dst, err := os.Create(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// update db data
	result, err := UpdateDoc(id,
		FilesData{
			Size:      handler.Size,
			CreatedAt: time.Now(),
			Url:       r.Host + "/assets/" + id.(primitive.ObjectID).Hex(),
			FilePath:  fp,
			Type:      handler.Header.Values("Content-Type")[0],
		})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		Del(id.(primitive.ObjectID))
		return
	}

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		Del(id.(primitive.ObjectID))
		return
	}

	// return json
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
