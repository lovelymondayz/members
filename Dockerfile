FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o main ./backend/cmd/server

# ─── Runtime ─────────────────────────────────────────
FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/main .

COPY .env.production .env

EXPOSE 8082

CMD ["./main"]
