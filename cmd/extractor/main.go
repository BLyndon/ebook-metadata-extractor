package main

import (
	"ebook-metadata-extractor/config"
	"ebook-metadata-extractor/pkg/fileutil"
	"ebook-metadata-extractor/pkg/metadata"
)

func main() {
	cfg := config.LoadConfig()

	titles := fileutil.ReadTitles(cfg)

	for _, title := range titles {
		metadata.ExtractMetaData(title, cfg)
	}
}
