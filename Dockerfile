FROM golang:1.22 AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -a -installsuffix cgo -o /extractor ./cmd/extractor

FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /

COPY --from=builder /extractor /extractor
COPY --from=builder /app/assets /assets

CMD ["/extractor"]
