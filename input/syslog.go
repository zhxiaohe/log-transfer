package input

import (
	"fmt"
	"log"
	"log_transfer/config"

	// "github.com/RichardKnop/machinery/v1/log"
	"gopkg.in/mcuadros/go-syslog.v2"
)

type syslogComponent struct {
	address      string
	protocol     string
	params       []string
	dataLocation string
}

func syslogInit() *syslogComponent {
	return &syslogComponent{
		address:  config.C.In.Syslog.Host,
		protocol: config.C.In.Syslog.Protocol,
	}
}

func (rs *syslogComponent) start(event chan map[string]interface{}) error {
	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)

	server := syslog.NewServer()

	// server.SetFormat(syslog.RFC5424)
	server.SetFormat(syslog.RFC3164)
	server.SetHandler(handler)

	if rs.protocol == "UDP" {
		server.ListenUDP(rs.address)
	} else {
		server.ListenUDP("0.0.0.0:514")
	}
	// else if rs.protocol == "TCP" { server.ListenTCP("0.0.0.0:7070")}
	log.Println("rsyslog server start!", rs.address)
	server.Boot()

	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			m := map[string]interface{}{}
			tag := fmt.Sprintf("%s", logParts["tag"])
			// fmt.Println(logParts["content"])
			// fmt.Println(logParts)
			logs := fmt.Sprintf("%s %s", tag, logParts["content"])
			m["value"] = logs
			m["key"] = tag
			select {
			case event <- m:
			default:
				log.Fatal("event channel is full!")
			}
		}
	}(channel)
	server.Wait()
	return nil
}
