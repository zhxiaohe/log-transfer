package output

import (
	"log"
	"log_transfer/config"
	"os"
	"sync"
)

type server interface {
	Start(chan map[string]interface{}) error
}

var handle server

func Start(event chan map[string]interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	switch config.C.Out.Type {
	case "kafka":
		handle = KafkaInit()
	case "stdout":
		handle = stdInit()
	default:
		log.Fatal("pls reset output server")
		os.Exit(1)
	}
	handle.Start(event)
}
