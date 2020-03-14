package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

const (
	PPROF_ADDR = ":9998"
)

func pprofInit() {
	pprofMux := http.DefaultServeMux
	go func() {
		server := &http.Server{
			Addr:    PPROF_ADDR,
			Handler: pprofMux,
		}
		log.Printf("[db] serving pprof on %s", PPROF_ADDR)
		panic(server.ListenAndServe())
	}()
	http.DefaultServeMux = http.NewServeMux()
}
