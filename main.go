package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type bufferedResponseWriter struct {
	delegate http.ResponseWriter
	status   int
	buffer   bytes.Buffer
	encoding string
}

func (brw *bufferedResponseWriter) Header() http.Header {
	return brw.delegate.Header()
}

func (brw *bufferedResponseWriter) Write(b []byte) (int, error) {
	brw.buffer.Write(b)
	return brw.delegate.Write(b)
}

func (brw *bufferedResponseWriter) WriteHeader(status int) {
	brw.status = status
	brw.delegate.WriteHeader(status)
}

func (brw *bufferedResponseWriter) String() string {
	var ret bytes.Buffer
	brw.delegate.Header().Write(&ret)
	ret.Write([]byte("\n"))
	switch brw.encoding {
	case "strip":
		t := transform.Chain(norm.NFKD, transform.RemoveFunc(func(r rune) bool {
			return r < 32 || r >= 127
		}))
		input := brw.buffer.Bytes()
		output := make([]byte, len(input))
		n, _, _ := t.Transform(output, input, true)
		ret.Write(output[:n])
	case "base64":
		ret.Write([]byte(base64.StdEncoding.EncodeToString(brw.buffer.Bytes())))
	case "none":
		ret.Write(brw.buffer.Bytes())
	}
	return ret.String()
}

func makeLogAndForwardHandler(proxy *httputil.ReverseProxy, encoding string) func(http.ResponseWriter, *http.Request) {
	log.Println("Created handler")
	var count int64
	var mutex = &sync.Mutex{}
	return func(w http.ResponseWriter, r *http.Request) {
		id := atomic.AddInt64(&count, 1)
		start := time.Now()
		req_dump, _ := httputil.DumpRequest(r, true)
		mutex.Lock()
		fmt.Printf("\n==== >>>REQUEST>>> %d @%s\n", id, start)
		fmt.Println(string(req_dump))
		mutex.Unlock()
		brw := bufferedResponseWriter{delegate: w, encoding: encoding}
		proxy.ServeHTTP(&brw, r)
		mutex.Lock()
		fmt.Printf("\n==== <<<RESPONSE<<< %d @%s [%s]\n", id, time.Now(), time.Since(start))
		fmt.Println(brw.String())
		fmt.Println()
		mutex.Unlock()
	}
}

const encodingHelp = `Output encoding.  One of none, base64, or strip.  None does no encoding
(dangerous if not redirecting output to a file as terminals may become confused
when images are printed).  Base64 base64 encodes responsees.  Strip removes
characters that may confuse terminals.`

func main() {
	log.SetOutput(os.Stdout)
	listenPort := flag.String("listenPort", "9999", "Listen Port")
	target := flag.String("target", "http://localhost:8500", "Target root url")
	encoding := flag.String("encoding", "strip", encodingHelp)
	insecure := flag.Bool("insecure", false, "Skip certificate verification when using SSL")
	flag.Parse()
	log.Println("Listening on port ", *listenPort)
	log.Println("Forwarding to ", *target)
	url, err := url.Parse(*target)
	if err != nil {
		log.Fatal(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	if *insecure {
		proxy.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	}
	http.HandleFunc("/", makeLogAndForwardHandler(proxy, *encoding))
	err = http.ListenAndServe(":"+*listenPort, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
