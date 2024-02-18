package main

import (
	"ebook-metadata-extractor/config"
	"log"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/metadata/extract", ExtractMetaDataHandler)

	port := config.LoadConfig().Port

	log.Println("Starting server on :" + strconv.Itoa(port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
