FROM golang:1.23-alpine as builder

RUN apk add --no-cache make g++

WORKDIR /app

COPY . .
RUN go mod download

RUN make build


FROM alpine:3.21 as relayer
WORKDIR /app

COPY --from=builder /app/bin/relayer .
COPY --from=builder /app/config.yaml .

ENTRYPOINT ["./relayer", "run"]
CMD ["--config", "config.yaml"]


FROM alpine:3.21 as resolver
WORKDIR /app

COPY --from=builder /app/bin/resolver .

ENTRYPOINT ["./resolver", "run"]
