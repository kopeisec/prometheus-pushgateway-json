FROM golang:1.26.1 AS builder
ENV CGO_ENABLED=0
WORKDIR /src
COPY . .
RUN go build -o ./bin/prometheus-pushgateway-json -mod vendor -trimpath -ldflags '-s -w' .

FROM debian:13
COPY --from=builder /src/bin/prometheus-pushgateway-json /usr/bin/prometheus-pushgateway-json
CMD ["/usr/bin/prometheus-pushgateway-json"]