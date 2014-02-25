package commands

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
	"github.com/roelrymenants/fileproxy/proxyconfig"
)

type RunCommand struct {
	IsVerbose bool
}

func ParseRunCommand(args []string) (Command, error) {
	runCmd := RunCommand{}

	runFlags := flag.NewFlagSet("run", flag.ExitOnError)

	runFlags.BoolVar(&runCmd.IsVerbose, "verbose", false, "Use -verbose to make the proxy verbose")
	runFlags.BoolVar(&runCmd.IsVerbose, "v", false, "Use -v to make the proxy verbose")
	
	runFlags.Parse(args)

	return &runCmd, nil
}

func (runCommand *RunCommand) Execute(configLoader ConfigLoader) error {
	config, err := configLoader()

	if err != nil {
		return err
	}

	proxy := goproxy.NewProxyHttpServer()

	//Verbose if either config or RunCommand is set to verbose
	proxy.Verbose = config.Verbose || runCommand.IsVerbose

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
			proxy.Verbose = config.Verbose
		}
	}()

	log.Printf("Starting proxyserver localhost:8080 with config %+v", *config)
	log.Fatal(http.ListenAndServe(":8080", proxy))

	return nil
}
