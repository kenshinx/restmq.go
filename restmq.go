package main

import (
	"fmt"
	"restmq.go/restmq"
)

var settings = restmq.Settings

func main() {
	fmt.Printf("RedisMQ :%s\n", settings.Version)
	httpServer := &restmq.HTTPServer{settings.HTTPServer}
	httpServer.Run()

}
