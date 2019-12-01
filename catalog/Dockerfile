FROM golang:1.13-alpine AS builder
RUN apk add --no-cache ca-certificates git
ENV PROJECT github.com/tony-yang/gcp-cloud-native-stack/catalog
WORKDIR /go/src/$PROJECT

COPY . .
RUN go build -o /catalog .

FROM alpine AS release
RUN apk add --no-cache ca-certificates \
 && GRPC_HEALTH_PROBE_VERSION=v0.3.1 \
 && wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 \
 && chmod +x /bin/grpc_health_probe
WORKDIR /catalog
COPY --from=builder /catalog ./server
COPY products.json .
EXPOSE 12000
ENTRYPOINT ["/catalog/server"]
