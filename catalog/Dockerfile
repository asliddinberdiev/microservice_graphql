FROM golang:1.23.3-alpine3.20 AS builder
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/asliddinberdiev/microservice_graphql
COPY go.mod go.sum ./
RUN go mod download
COPY catalog catalog
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/app ./catalog/cmd/catalog

FROM alpine:3.20
WORKDIR /usr/bin
COPY --from=builder /go/bin .
EXPOSE 8080
CMD ["./app"]