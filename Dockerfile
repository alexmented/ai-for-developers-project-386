FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY backend ./backend

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/server ./backend/cmd/server

FROM alpine:3.21

RUN apk add --no-cache curl

WORKDIR /app

COPY --from=builder /out/server /app/server

EXPOSE 10000

HEALTHCHECK --interval=5s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:${PORT:-10000}/public/name-owner || exit 1

CMD ["/app/server"]
