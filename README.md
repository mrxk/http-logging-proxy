# http-logging-proxy

HTTP Logging Proxy

## Description

This project builds a simple reverse proxy that logs requests and responses to
stdout.  It can be useful when debugging/investigating REST API trafic.  It is
essentially a wrapper around the golang
[SingleHostReverseProxy](https://pkg.go.dev/net/http/httputil#NewSingleHostReverseProxy).

## Usage
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
