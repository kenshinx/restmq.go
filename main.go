package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/kenshinx/restmq.go/restmq"
	"os"
)

var (
	Settings   restmq.RestMQSettings
	configFile string
	protocol   string
)

func main() {
	fmt.Printf("RedisMQ :%s\n", Settings.Version)

	restmq.Settings = Settings
	httpServer := restmq.HTTPServer{}
	go httpServer.Run(protocol)

	<-make(chan bool)

}

func init() {

	var configFile string

	flag.StringVar(&configFile, "c", "restmq.conf", "Look for restmq toml-formatting config file in this directory")
	flag.StringVar(&protocol, "p", "http", "The webserver protocol,usage: {http|fcgi|scgi}")
	flag.Parse()

	if _, err := toml.DecodeFile(configFile, &Settings); err != nil {
		fmt.Printf("%s is not a valid toml config file\n", configFile)
		fmt.Println(err)
		os.Exit(1)
	}

}
