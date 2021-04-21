FROM golang:alpine AS builder

RUN apk --no-cache add git gcc

ADD . /go
RUN cd /go && go build -o goc

FROM golang:alpine

RUN apk --no-cache add git gcc

COPY --from=builder /go/goc /bin/goc