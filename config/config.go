package config

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Mysql Mysql `yaml:"mysql"`
	Redis Redis `yaml:"redis"`
	In    In    `yaml:"in"`
	Out   Out   `yaml:"out"`
	Event Event `yaml:"event"`
	Pprof Pprof `yaml:"pprof"`
}

type Mysql struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type Redis struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password""`
	Database string `yaml:"database""`
}

type In struct {
	Type   string `yaml:"type"`
	Syslog struct {
		Host     string `yaml:"host"`
		Protocol string `yaml:"protocol"`
	} `yaml:"syslog"`
	HTTP struct {
		Addr string `yaml:"addr"`
	} `yaml:"http"`
	Grpc struct {
		Addr string `yaml:"addr"`
	} `yaml:"grpc"`
}

type Out struct {
	Type  string `yaml:"type"`
	Kafka struct {
		Brokers  string `yaml:"brokers"`
		Version  string `yaml:"version"`
		Group    string `yaml:"group"`
		Topics   string `yaml:"topics"`
		Assignor string `yaml:"assignor"`
		Oldest   bool   `yaml:"oldest"`
		Verbose  bool   `yaml:"verbose"`
	} `yaml:"kafka"`
}

type Event struct {
	Chansize int `yaml:"chansize"`
}
type Pprof struct {
	Enable bool   `yaml:"enable"`
	Addr   string `yaml:"addr"`
}

var C = Config{}

func Parse(c string) {
	f, err := os.Open(c)
	if err != nil {
		pwd, _ := os.Getwd()
		log.Fatal(err, " ", pwd)
	}
	content, err := ioutil.ReadAll(f)
	err = yaml.Unmarshal([]byte(content), &C)
	// fmt.Println(C)
	if C.Event.Chansize == 0 {
		C.Event.Chansize = 20000
	}
	defer f.Close()
}
