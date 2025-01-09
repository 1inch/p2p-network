FROM golang:1.23-alpine as builder

RUN apk add --no-cache make

WORKDIR /app

COPY . .

RUN go mod download

ENTRYPOINT ["make", "build_resolver"]


FROM alpine:3.21 as resolver
WORKDIR /app

COPY --from=builder /app/bin/resolver .
COPY --from=builder /app/assets/resolver.config.yaml /app/config.yaml

ENTRYPOINT ["./resolver", "run"]
CMD ["--configFile", "config.yaml"]
