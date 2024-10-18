FROM golang:1.21
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./src/ /app/src
COPY ./fonts/ /app/fonts
COPY config.json ./
RUN CGO_ENABLED=0 GOOS=linux go build -o ./tml-readme-card ./src

EXPOSE 8005
CMD ["./tml-readme-card"]
