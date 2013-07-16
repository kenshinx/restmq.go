package restmq

import (
	"strconv"
)

type RestMQSettings struct {
	Version    string
	Debug      bool
	HTTPServer WebServerSettings `toml:"http"`
	Redis      RedisSettings     `toml:"redis"`
	Log        LogSettings       `toml:"log"`
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

type LogSettings struct {
	File string
}
