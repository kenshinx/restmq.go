package main

import (
	"fmt"
	"github.com/kenshinx/restmq.go/restmq"
)

var settings = restmq.Settings

func main() {
	fmt.Printf("RedisMQ :%s\n", settings.Version)

	httpServer := restmq.HTTPServer{}
	go httpServer.Run()

	<-make(chan bool)

}
