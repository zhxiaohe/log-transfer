package input

import (
	"log"
	"log_transfer/config"
	"sync"
)

type input interface {
	start(chan map[string]interface{}) error
}

var handle input

func Start(event chan map[string]interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	switch config.C.In.Type {
	case "syslog":
		handle = syslogInit()
	case "http":
		handle = httpInit()
	case "beat":
		log.Fatal("pls reset input server")
	case "grpc":
		handle = grpcInit()
	default:
		log.Fatal("pls reset input server")
	}
	handle.start(event)
	defer log.Panicln("exit!!!")
}
