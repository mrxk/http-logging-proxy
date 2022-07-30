FROM golang:1.18-buster AS builder

WORKDIR /go/http-logging-proxy

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -o executable

FROM scratch

COPY --from=builder /go/http-logging-proxy/executable /executable

ENTRYPOINT ["/executable"]
