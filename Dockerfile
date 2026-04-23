FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git make

COPY go.mod go.sum* ./
RUN go mod download || true

COPY . .
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/server

# Ensure web/public exists so the COPY command in the next stage doesn't fail
RUN mkdir -p /app/web/public

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

# Run as non-root user for security
RUN addgroup -S edusys && adduser -S edusys -G edusys

WORKDIR /app

COPY --from=builder /server .
COPY --from=builder /app/web/public ./web/public

# Do NOT copy .env.example as .env — secrets must be injected via env vars at runtime

RUN chown -R edusys:edusys /app

USER edusys

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget -qO- http://localhost:8080/ || exit 1

CMD ["./server"]