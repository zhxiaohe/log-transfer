package main

import (
	"flag"
	"fmt"
	"log_transfer/config"
	"log_transfer/input"
	"log_transfer/output"
	"log_transfer/status"
	"sync"
)

var conf string

func init() {
	flag.StringVar(&conf, "c", "", "config file path")
	flag.Parse()

}

func main() {
	config.Parse(conf)
	Event := make(chan map[string]interface{}, config.C.Event.Chansize)
	var wg sync.WaitGroup
	wg.Add(3)
	fmt.Println(conf)
	go input.Start(Event, &wg)
	go output.Start(Event, &wg)
	go status.Start(&wg)
	wg.Wait()
}
