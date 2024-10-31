package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"xgateway/internal/conf"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

var flagconf string

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
	var cstZone = time.FixedZone("CST", 8*3600) // 东八
	time.Local = cstZone
}

func Print(handler runtime.HandlerFunc) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		log.Printf("received request: %+v", r)
		log.Printf("received params: %+v", params)
		handler(w, r, params)
	}
}

func newGwMux() *runtime.ServeMux {
	return runtime.NewServeMux(
		runtime.WithMiddlewares(Print),
	)
}

func main() {
	flag.Parse()
	c, err := conf.LoadConfig(flagconf)
	if err != nil {
		log.Fatalf("load config error: %v", err)
	}
	s, err := wireApp(c)
	if err != nil {
		log.Fatalf("wire app error: %v", err)
	}
	s.Start()
	defer s.Stop()
	gwServer := &http.Server{
		Addr:    c.Gateway.Addr,
		Handler: s.gw,
	}
	log.Printf("start gateway server at %s", c.Gateway.Addr)
	err = gwServer.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to start gateway server: %v", err)
	}
}
