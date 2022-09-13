package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	// Maximum upload of 10 MB files
	r.ParseMultipartForm(10 << 20)

	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	var assetsPath = "public/assets/" + time.Now().Format("2006-01-02")

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

	fp := assetsPath + "/" + id.(primitive.ObjectID).Hex() + filepath.Ext(handler.Filename)

	// Create file
	dst, err := os.Create(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	result, err := UpdateDoc(id, FilesData{Size: handler.Size, CreatedAt: time.Now(), Url: r.Host + "/assets/" + id.(primitive.ObjectID).Hex(), FilePath: fp})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
	// fmt.Fprintf(w, "Successfully Uploaded File\n")
}
