package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
	"github.com/roelrymenants/fileproxy"
)

func main() {
	config := *fileproxy.LoadConfig(fileproxy.DefaultConfigFile)

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	var matchCondition = goproxy.ReqConditionFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		_, ok := config.Rewrites[req.URL.Path]

		log.Printf("Trying to match %s", req.URL.Path)

		if ok {
			log.Printf("Matched %s", req.URL.Path)
		}

		return ok
	})

	proxy.OnRequest(matchCondition).DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		dstFile, _ := config.Rewrites[req.URL.Path]
		fileBytes, err := ioutil.ReadFile(dstFile)
		if err == nil {
			return req, goproxy.NewResponse(req, "text/css", 200, string(fileBytes[:]))
		} else {
			log.Printf("Error reading file %s", dstFile)
			return req, nil
		}
	})

	configChan := fileproxy.StartWatching(fileproxy.DefaultConfigFile)

	go func() {
		for {
			config = <-configChan
		}
	}()

	log.Fatal(http.ListenAndServe(":8080", proxy))
}
