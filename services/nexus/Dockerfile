FROM golang:1.27.1-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o gristguard .

FROM alpine:3
WORKDIR /app

RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S -G app app

COPY --from=builder /app/gristguard .

USER app
EXPOSE 8080

ENTRYPOINT ["./gristguard"]
