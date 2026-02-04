package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/storage"
)

func main() {
	http.HandleFunc("/", handleRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if ok := r.Header.Get("x-goog-iap-jwt-assertion"); ok == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/")

	if path == "" {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Usage:\n  GET /gcs_bucket/path/to/object"))
		return
	}

	bucketName, objectName, found := strings.Cut(path, "/")
	if !found {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)
	object := bucket.Object(objectName)

	// オブジェクトを読み取ってレスポンスに流す
	reader, err := object.NewReader(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			http.Error(w, "File not found", http.StatusNotFound)
		} else if os.IsPermission(err) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		} else {
			log.Printf("Failed to read object: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	defer reader.Close()

	// Content-Type を設定
	if reader.Attrs.ContentType != "" {
		w.Header().Set("Content-Type", reader.Attrs.ContentType)
	}

	// コンテンツをストリーム
	if _, err := io.Copy(w, reader); err != nil {
		log.Printf("Failed to stream object: %v", err)
	}
}
