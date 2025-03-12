# Build stage
FROM golang:1.22.4 AS builder

WORKDIR /app

COPY go.mod .
COPY server-single/main.go .
COPY cert.crt .
COPY cert.key .

RUN go build -o server-single main.go

# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server-single .
COPY cert.crt .
COPY cert.key .

EXPOSE 6660

CMD ["./server-single"]
