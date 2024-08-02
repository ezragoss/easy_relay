FROM golang:1.22.5-alpine3.20 as builder

WORKDIR /

COPY src .

RUN go mod download

RUN go build -o /relay

FROM alpine:3.20

WORKDIR /

COPY --from=builder /relay /bin/relay

EXPOSE 80

ENTRYPOINT ["relay"]