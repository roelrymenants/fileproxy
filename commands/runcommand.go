package commands

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
	"github.com/roelrymenants/fileproxy/proxyconfig"
)

type RunCommand struct {
}

func ParseRunCommand(_ []string) (*RunCommand, error) {
	return &RunCommand{}, nil
}

func (runCommand *RunCommand) Execute(config *proxyconfig.Config) error {
	proxy := goproxy.NewProxyHttpServer()

	proxy.Verbose = config.Verbose

	var matchCondition = goproxy.ReqConditionFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		_, ok := config.Rewrites[req.URL.String()]

		log.Printf("Trying to match %s", req.URL.String())

		if ok {
			log.Printf("Matched %s", req.URL.String())
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

	configChan := proxyconfig.StartWatching(proxyconfig.DefaultConfigFile)

	go func() {
		for {
			config = <-configChan
		}
	}()

	log.Fatal(http.ListenAndServe(":8080", proxy))

	return nil
}
