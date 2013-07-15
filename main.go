package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"restmq.go/restmq"
)

var (
	Settings restmq.RestMQSettings
)

func main() {
	fmt.Printf("RedisMQ :%s\n", Settings.Version)

	restmq.Settings = Settings
	httpServer := restmq.HTTPServer{}
	go httpServer.Run()

	<-make(chan bool)

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
