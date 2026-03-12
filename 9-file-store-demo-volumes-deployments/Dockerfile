# Build stage
FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod ./
COPY main.go ./

RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o gofile-store main.go

# Runtime stage
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/gofile-store /app/gofile-store

EXPOSE 8080

ENTRYPOINT ["/app/gofile-store"]