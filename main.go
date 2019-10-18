package main

import (
	"gopkg.in/elazarl/goproxy.v1"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			r.Header.Set("X-GoProxy", "shaswata56")
			return r, nil
		})
	log.Fatal(http.ListenAndServe(":"+port, proxy))
}
