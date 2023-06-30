FROM golang:1.20.5-buster
WORKDIR /app
COPY example-golang.go ./

RUN go build -o /usr/bin/go-example ./example-golang.go
