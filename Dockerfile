FROM golang:1.23.8-bullseye AS builder
WORKDIR /src
COPY . .
RUN go build -o ./bin/prometheus-pushgateway-json -mod vendor -trimpath -ldflags '-s -w' .

FROM debian:11
COPY --from=builder /src/bin/prometheus-pushgateway-json /usr/bin/prometheus-pushgateway-json
CMD ["/usr/bin/prometheus-pushgateway-json"]