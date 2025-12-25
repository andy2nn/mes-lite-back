# ---------- build ----------
FROM golang:1.25 AS builder

WORKDIR /app

# зависимости
COPY go.mod go.sum ./
RUN go mod download

# исходники
COPY . .

# сборка
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o api ./cmd/server

# ---------- runtime ----------
FROM alpine:3.19

WORKDIR /app

# сертификаты для HTTPS / JWT / PostgreSQL
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/api .

EXPOSE 8080

CMD ["./api"]
