# syntax=docker/dockerfile:1
FROM golang:1.22.5-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o statistic-consumer ./main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/statistic-consumer .
EXPOSE 8080
CMD ["./statistic-consumer"]
