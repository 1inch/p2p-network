FROM golang:1.23-alpine as builder

WORKDIR /app

COPY ./../ .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o resolver ./cmd/resolver/main.go

FROM alpine:3.21 as resolver
WORKDIR /app

COPY --from=builder /app/resolver .
COPY --from=builder /app/assets/resolver.config.yaml /app/config.yaml

RUN apk update && apk upgrade --no-cache libcrypto3 libssl3
RUN adduser -u 1000 resolver --disabled-password
USER resolver