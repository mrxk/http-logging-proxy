FROM golang:1.17-buster AS builder

WORKDIR /go/http-logging-proxy

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -o http-logging-proxy

FROM scratch

COPY --from=builder /go/http-logging-proxy/http-logging-proxy /http-logging-proxy

ENTRYPOINT ["/http-logging-proxy"]
