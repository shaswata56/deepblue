package main

import (
	"flag"
	"gopkg.in/elazarl/goproxy.v1"
	"log"
	"net/http"
	"os"
	"regexp"
)

func main() {
	port := os.Getenv("PORT")
	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("^.*$"))).
		HandleConnect(goproxy.AlwaysMitm)
	verbose := flag.Bool("v", false, "should every proxy request be logged to stdout")
	addr := flag.String("", ":"+port, "proxy listening")
	flag.Parse()
	proxy.Verbose = *verbose
	log.Fatal(http.ListenAndServe(*addr, proxy))
}
