# Build stage
FROM golang:1.27.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/qrgen ./cmd

# Runtime stage
FROM alpine:3.22

WORKDIR /app
RUN apk --no-cache add ca-certificates

COPY --from=builder /out/qrgen ./qrgen
COPY config/ ./config/

EXPOSE 8080

CMD ["./qrgen"]
