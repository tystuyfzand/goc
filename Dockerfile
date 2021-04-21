FROM golang:alpine

RUN apk --no-cache add git gcc

RUN go build -o /bin/goc