FROM golang:alpine AS builder

RUN apk --no-cache add git gcc

ADD . /app
RUN cd /app && go build -o goc

FROM golang:alpine

RUN apk --no-cache add git gcc

COPY --from=builder /app/goc /bin/goc