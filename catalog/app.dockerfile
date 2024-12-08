FROM golang:1.23.3-alpine3.20 AS builder
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/asliddinberdiev/microservice_graphql
COPY go.mod go.sum ./
COPY vendor vendor
COPY catalog catalog
RUN GO111MODULE=on go build -mod vendor -o /go/bin/catalog ./catalog/cmd/catalog

FROM alpine:3.20
WORKDIR /usr/bin
COPY --from=builder /go/bin .
EXPOSE 8080
CMD ["/catalog"]