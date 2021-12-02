package status

import (
	"log"
	"log_transfer/config"
	"net/http"
	_ "net/http/pprof"
	"sync"
)

func Start(wg *sync.WaitGroup) {
	defer wg.Done()
	if !config.C.Pprof.Enable {
		return
	}
	// http.HandleFunc("status", handler func(http.ResponseWriter, *http.Request))
	// http.HandleFunc("/status/pprof/", pprof.Index)
	// http.HandleFunc("/status/pprof/cmdline", pprof.Cmdline)
	// http.HandleFunc("/status/pprof/profile", pprof.Profile)
	// http.HandleFunc("/status/pprof/symbol", pprof.Symbol)
	// http.HandleFunc("/status/pprof/trace", pprof.Trace)
	log.Println("DEBUG pport start", config.C.Pprof.Addr)
	http.ListenAndServe(config.C.Pprof.Addr, nil)
}
