package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"encoding/json"

	"github.com/elazarl/goproxy"
)

type Config struct {
	Rewrites map[string]string
}

func LoadConfig(filepath string) Config {
	var config Config
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Panic("Error reading config file", err)
	}
	err = json.Unmarshal(b, &config)
	log.Printf("%v", config)

	return config
}

func main() {
	config := LoadConfig("rewrites.json")

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

	log.Fatal(http.ListenAndServe(":8080", proxy))
}
