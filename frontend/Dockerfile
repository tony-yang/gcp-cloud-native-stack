FROM golang:1.13-alpine AS builder
RUN apk add --no-cache ca-certificates git
ENV PROJECT github.com/tony-yang/gcp-cloud-native-stack/frontend
WORKDIR /go/src/$PROJECT

COPY . .
RUN go build -o /frontend .

FROM alpine AS release
RUN apk add --no-cache \
    bind-tools \
    busybox-extras \
    ca-certificates \
    net-tools
WORKDIR /frontend
COPY --from=builder /frontend ./server
COPY ./templates ./templates
EXPOSE 13000
ENTRYPOINT ["/frontend/server"]
