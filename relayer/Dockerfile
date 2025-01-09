FROM golang:1.23-alpine as builder

RUN apk add --no-cache make

WORKDIR /app

COPY . .

RUN go mod download

ENTRYPOINT ["make", "build_relayer"]


FROM alpine:3.21 as relayer
WORKDIR /app

COPY --from=builder /app/bin/relayer .
COPY --from=builder /app/assets/relayer.config.yaml /app/config.yaml

ENTRYPOINT ["./relayer", "run"]
CMD ["--config", "config.yaml"]
