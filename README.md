# http-logging-proxy

HTTP Logging Proxy

## Description

This project builds a simple reverse proxy that logs requests and responses to
stdout.  It can be useful when debugging/investigating REST API trafic.  It is
essentially a wrapper around the golang
[SingleHostReverseProxy](https://pkg.go.dev/net/http/httputil#NewSingleHostReverseProxy).

## CLI Usage
```
Usage of http-logging-proxy:
  -encoding string
    	Output encoding.  One of none, base64, or strip.  None does no encoding
    	        (dangerous if not redirecting output to a file as terminals may become
    	        confused when images are printed).  Base64 base64 encodes responsees.  Strip
    	        removes characters that may confuse terminals.  Default: strip (default "strip")
  -listenPort string
    	Listen Port (default "9999")
  -target string
    	Target root url (default "http://localhost:8500")
```

## Conainterized Usage
Build the image:
```
ubuntu@ip-172-31-28-115:~/code/http-logging-proxy$ make build
docker build -t http-logging-proxy .
Sending build context to Docker daemon  73.22kB
Step 1/8 : FROM golang:1.18-buster AS builder
 ---> 0e87973a8632
Step 2/8 : WORKDIR /go/http-logging-proxy
 ---> Using cache
 ---> baa27d2256fb
Step 3/8 : COPY . .
 ---> 97fea7b208a4
Step 4/8 : RUN go mod download
 ---> Running in 1397869b8c31
Removing intermediate container 1397869b8c31
 ---> a4612bf82823
Step 5/8 : RUN CGO_ENABLED=0 go build -o executable
 ---> Running in 0560e76753c9
Removing intermediate container 0560e76753c9
 ---> eb9d34ed2720
Step 6/8 : FROM scratch
 ---> 
Step 7/8 : COPY --from=builder /go/http-logging-proxy/executable /executable
 ---> Using cache
 ---> 71533a3b8812
Step 8/8 : ENTRYPOINT ["/executable"]
 ---> Using cache
 ---> 507715478a6d
Successfully built 507715478a6d
Successfully tagged http-logging-proxy:latest
```

Use with other containers, such as in Docker Compose:
```
  http-logging-proxy:
    image: http-logging-proxy:latest
    command:
      # listen on port 80 for incoming traffic and pass it to a server running on port 5000 in the same compose file
      -listenPort 80 -target http://server:5000
    ports:
      # expose this endpoint on port 8080
      - 8080:80
```
