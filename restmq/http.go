package restmq

import (
	"encoding/json"
	"fmt"
	"github.com/hoisie/redis"
	"github.com/hoisie/web"
	"github.com/kenshinx/redisq"
	"log"
	"os"
)

type IndexHandler struct {
}

func (h *IndexHandler) Get() string {
	return "Hello RestMQ !"
}

type RestQueueHandler struct {
	redis  *redis.Client
	logger *log.Logger
}

func (h *RestQueueHandler) Queue(name string) (queue *redisq.RedisQueue) {
	queue = redisq.NewRedisQueue(h.redis, name)
	queue.SetLogger(h.logger)
	return
}

func (h *RestQueueHandler) Get(ctx *web.Context, val string) {
	queue := h.Queue(val)
	if queue.Empty() {
		ctx.NotFound(Status.QueueNotFound(val))
		return
	}
	msg, err := queue.GetNoWait()
	resp, err := json.Marshal(msg)
	//msg.(type) is iteface{}
	//resp.(type) is []byte
	if err != nil {
		ctx.ResponseWriter.WriteHeader(400)
		ctx.WriteString(Status.BadRequest(val))
		if Settings.Debug {
			debug := fmt.Sprintf("Debug: %s", err)
			ctx.WriteString(debug)
		}
		h.logger.Fatalf("Dequeue from <%s> Error:%s", val, err)
		return
	}
	ctx.SetHeader("Content-Type", "application/json; charset=UTF-8", true)
	ctx.WriteString(string(resp))

}

func (h *RestQueueHandler) Put(ctx *web.Context, val string) {
	h.logger.Println(ctx)
	h.logger.Println(val)
}

func (h *RestQueueHandler) Clear(ctx *web.Context, val string) {
	h.logger.Println(ctx)
	h.logger.Println(val)
}

func initLogger(log_file string) (logger *log.Logger) {
	if log_file != "" {
		f, err := os.Create(log_file)
		if err != nil {
			os.Exit(1)
		}
		logger = log.New(f, "[http-webserver]", log.Ldate|log.Ltime)
	} else {
		logger = log.New(os.Stdout, "[http-webserver]", log.Ldate|log.Ltime)
	}
	return logger

}

type HTTPServer struct {
}

func (s HTTPServer) Run() {

	logger := initLogger(Settings.Log.File)
	redis := &redis.Client{Addr: Settings.Redis.Addr(),
		Db:       Settings.Redis.DB,
		Password: Settings.Redis.Password}

	var (
		indexHandler = &IndexHandler{}
		queueHandler = &RestQueueHandler{redis, logger}
	)

	web.Get("/", indexHandler.Get)
	web.Get("/q/(.+)", queueHandler.Get)
	web.Post("/q/(.+)", queueHandler.Put)
	web.Delete("/q/(.+)", queueHandler.Clear)
	web.SetLogger(logger)
	web.Run(Settings.HTTPServer.Addr())
}
