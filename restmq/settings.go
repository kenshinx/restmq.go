package restmq

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"strconv"
)

var (
	Settings RestMQSettings
)

func Debug(format string, items ...interface{}) {
	if Settings.Debug {
		log.Printf(format, items...)
	}
}

type RestMQSettings struct {
	Version    string
	Debug      bool
	HTTPServer WebServerSettings `toml:"http"`
	WSServer   WebServerSettings `toml:websocket`
	Redis      RedisSettings     `toml:"redis"`
}

type WebServerSettings struct {
	Host string
	Port int
}

func (s *WebServerSettings) Addr() string {
	return s.Host + ":" + strconv.Itoa(s.Port)
}

type RedisSettings struct {
	Host     string
	Port     int
	DB       int
	Password string
}

func (s *RedisSettings) Addr() string {
	return s.Host + ":" + strconv.Itoa(s.Port)
}

func init() {

	var configFile string

	flag.StringVar(&configFile, "c", "restmq.conf", "Look for restmq toml-formatting config file in this directory")
	flag.Parse()

	if _, err := toml.DecodeFile(configFile, &Settings); err != nil {
		fmt.Printf("%s is not a valid toml config file\n", configFile)
		fmt.Println(err)
		os.Exit(1)
	}

}
