package main

import (
	"ebook-metadata-extractor/config"
	"ebook-metadata-extractor/pkg/metadata"
	"fmt"
	"net/http"
)

func ExtractMetaDataHandler(w http.ResponseWriter, r *http.Request) {
	cfg := config.LoadConfig()

	title := r.URL.Query().Get("title")
	if title == "" {
		http.Error(w, "Missing title query parameter", http.StatusBadRequest)
		return
	}

	jsonResponse, err := metadata.ExtractMetaData(title, cfg)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error extracting metadata: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonResponse))
}
