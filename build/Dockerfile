FROM golang:1.21.5-alpine3.18 as builder

WORKDIR /fio

#RUN apk add build-base
#RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY ../go.mod go.sum ./
RUN go mod download

ADD .. /fio/

RUN go build ./cmd/fio

ENTRYPOINT ./fio
