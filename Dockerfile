ARG GO_VERSION=1.19

# Build with the latest Go image
FROM golang:alpine AS builder

RUN apk --no-cache add git gcc

ADD . /app
RUN cd /app && go build -o goc

# Run inside the specified Go image
FROM golang:${GO_VERSION}-alpine

RUN apk --no-cache add git gcc

COPY --from=builder /app/goc /bin/goc