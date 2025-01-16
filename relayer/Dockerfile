FROM golang:1.23-alpine as builder

WORKDIR /app

COPY ./../ .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o relayer ./cmd/relayer/main.go

FROM alpine:3.21 as relayer
WORKDIR /app

COPY --from=builder /app/relayer .
COPY --from=builder /app/assets/relayer.config.yaml /app/config.yaml

