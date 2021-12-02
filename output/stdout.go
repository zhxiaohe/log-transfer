package output

import (
	"fmt"
)

type stdoutComponent struct {
}

func stdInit() *stdoutComponent {
	return &stdoutComponent{}
}

func (s *stdoutComponent) Start(event chan map[string]interface{}) error {
	go func() {
		for {
			select {
			case m := <-event:
				fmt.Println(m["value"])
				// default:
				// log.Fatal("stdout exit")
				// return
			}
		}
	}()
	return nil
}
