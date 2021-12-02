package input

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log_transfer/config"
	"net/http"
)

type httpserver struct {
	event chan map[string]interface{}
}

type rbody struct {
	Tag  string
	Data string
}

func httpInit() *httpserver {
	return &httpserver{}
}

func (h *httpserver) start(event chan map[string]interface{}) error {
	h.event = event
	http.HandleFunc("/put", h.viewfunc)
	if config.C.In.HTTP.Addr == "" {
		http.ListenAndServe(":7808", nil)
	}
	http.ListenAndServe(config.C.In.HTTP.Addr, nil)

	return nil
}

func (h *httpserver) viewfunc(w http.ResponseWriter, r *http.Request) {
	//POST {"tag":"b", "data":"aaaa"}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	HTTPReturn := map[string]string{}
	err := r.ParseForm()
	if err != nil {
		HTTPReturn["status"] = "401"
		jsonData, _ := json.Marshal(HTTPReturn)
		io.WriteString(w, string(jsonData))
		return
	}
	body, _ := ioutil.ReadAll(r.Body)
	var d rbody
	json.Unmarshal([]byte(body), &d)
	m := map[string]interface{}{}
	m["value"] = d.Data
	m["tag"] = d.Tag
	select {
	case h.event <- m:
	}
	HTTPReturn["status"] = "200"
	jsonData, _ := json.Marshal(HTTPReturn)
	io.WriteString(w, string(jsonData))
}
