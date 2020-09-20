FROM golang:1-alpine

COPY . .

CMD go run cmd/main.go
