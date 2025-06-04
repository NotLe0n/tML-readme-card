FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./src/ /app/src
COPY ./fonts/ /app/fonts
COPY config.json ./

RUN CGO_ENABLED=0 GOOS=linux go build -o tml-readme-card ./src

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    libssl3 \
    ca-certificates \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/tml-readme-card .
COPY --from=builder /app/fonts ./fonts
COPY --from=builder /app/config.json .

EXPOSE 8005

CMD ["./tml-readme-card"]
