# ========== Build Stage ==========
FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /server ./cmd/server

# ========== Runtime Stage ==========
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata curl

RUN addgroup -S edusys && adduser -S edusys -G edusys

WORKDIR /app

COPY --from=builder /server .

RUN chown -R edusys:edusys /app

USER edusys

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

CMD ["./server"]