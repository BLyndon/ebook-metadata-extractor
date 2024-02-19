package main

import (
	"ebook-metadata-extractor/cmd/extractor/middleware"
	"ebook-metadata-extractor/config"
	"log"
	"net/http"
	"strconv"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/metadata/extract", ExtractMetaDataHandler)

	// Wrap the entire mux with the CORS middleware
	handler := middleware.Cors(mux)

	port := config.LoadConfig().Port
	log.Println("Starting server on :" + strconv.Itoa(port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), handler))
}
