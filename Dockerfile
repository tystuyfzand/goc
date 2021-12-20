FROM golang:alpine

RUN apk --no-cache add git gcc

COPY /build/goc_linux_amd64 /bin/goc