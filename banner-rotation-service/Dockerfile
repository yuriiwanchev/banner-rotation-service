# syntax=docker/dockerfile:1
FROM golang:1.18-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o banner-rotation-service ./cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/banner-rotation-service .
EXPOSE 8080
CMD ["./banner-rotation-service"]
