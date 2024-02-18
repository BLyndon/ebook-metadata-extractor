all: build, run

build:
	@go build -o ./bin/extractor ./cmd/extractor

run: build
	@./bin/extractor
