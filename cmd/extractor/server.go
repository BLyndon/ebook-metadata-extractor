package main

import (
	"ebook-metadata-extractor/config"
	"ebook-metadata-extractor/pkg/fileutil"
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

	metadataString, err := fileutil.ReadMetadataIfExists(title, cfg)
	if err != nil {
		metadataString, err = metadata.ExtractMetaData(title, cfg)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error retrieving metadata: %v", err), http.StatusInternalServerError)
			return
		}
	}

	if cfg.PeristMetadata {
		err := fileutil.WriteToFile(metadataString, title, cfg)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error persisting metadata: %v", err), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(metadataString))
}
