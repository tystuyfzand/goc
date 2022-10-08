ARG BASE_IMAGE=golang:alpine

# Build with the latest Go image
FROM golang:alpine AS builder

RUN apk --no-cache add git gcc

ADD . /app
RUN cd /app && go build -o goc

# Run inside the specified Go image
FROM $BASE_IMAGE

RUN apk --no-cache add git gcc

COPY --from=builder /app/goc /bin/goc