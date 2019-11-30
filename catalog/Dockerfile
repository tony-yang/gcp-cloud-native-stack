FROM golang:1.13-alpine AS builder
RUN apk add --no-cache ca-certificates git
ENV PROJECT github.com/tony-yang/gcp-cloud-native-stack/catalog
WORKDIR /go/src/$PROJECT

COPY . .
RUN go mod init \
 && go build -o /catalog .
