# eBook Metadata Extractor (eMdE)

## Overview

eMdE is a Go-based tool designed to automate the extraction of metadata from eBooks and organize them efficiently. You can query metadata for a certain eBook via a REST API. Than the tool checks if the metadata is already present in the database. If not, it extracts the metadata from the eBook title and returns it in a structured JSON format. Depending on the configuration, the metadata is persisted.

## Features

- Extracts eBook metadata using OpenAI API or loads it from a database.
- Outputs metadata in a structured JSON format following a predefined schema.
- Easy to configure and extend for additional metadata fields or eBook formats.

## Getting Started

### Prerequisites

- Go 1.22
- Access to OpenAI API and a valid API key.

### Configuration

1. Set your OpenAI API key as an environment variable:

   ```sh
   export OPENAI_API_KEY="your_openai_api_key_here"
   ```

2. Adjust the configuration in `config/config.go` as needed.

## Usage

There is a Makefile provided to build and run the application.

- Run locally:

  ```sh
  make run
  ```

- Run docker container:

  ```sh
   make docker-run
  ```

Chechout the Makefile for more commands.

## Metadata API

The API is a simple REST API that allows you to query metadata for a given title

```sh
curl -X GET "http://localhost:8080/metadata/extract?title=Wiley%20Finance%20John%20C%20Hull%20-%20Risk%20Management%20and%20Financial%20Institutions%202018%20Wiley" -H "accept: application/json"
```