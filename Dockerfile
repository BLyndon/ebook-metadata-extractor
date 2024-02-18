# Build stage
FROM golang:1.22 AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.* ./

# Download Go modules
RUN go mod download

# Copy the source code into the container
COPY . .

# Set environment variables for static build
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# Build the Go app
RUN go build -a -installsuffix cgo -o /extractor ./cmd/extractor

# Ensure the binary is executable
RUN chmod +x /extractor

# Final stage
FROM alpine:latest  

# Add CA certificates
RUN apk --no-cache add ca-certificates

WORKDIR /

# Copy the binary from the builder stage
COPY --from=builder /extractor /extractor

CMD ["/extractor"]
