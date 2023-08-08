FROM golang:1.18.2-alpine

WORKDIR /app
COPY . .
RUN go build -o app

ENTRYPOINT ["./app"]
