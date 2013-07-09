package restmq

import (
	"github.com/hoisie/web"
	"log"
)

type HTTPServer struct {
	Setting WebServerSettings
}

func (s *HTTPServer) Run() {
	var server web.Server
	log.Printf("HTTP server begin running on %s\n", s.Setting.Addr())
	server.Run(s.Setting.Addr())
}
