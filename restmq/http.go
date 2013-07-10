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

func (h *RestQueueHandler) List(ctx *web.Context) {

	var keys []string
	keys, _ = h.redis.Keys("*")
	resp, _ := json.Marshal(keys)
	ctx.SetHeader("Content-Type", "application/json; charset=UTF-8", true)
	ctx.WriteString(string(resp))
}

func (h *RestQueueHandler) Get(ctx *web.Context, val string) {
	queue := h.Queue(val)
	if !queue.Exists() {
		ctx.NotFound(QueueNotFound)
		return
	}
	if queue.Empty() {
		ctx.ResponseWriter.WriteHeader(400)
		ctx.WriteString(EmptyQueue)
		return
	}
	mesg, err := queue.GetNoWait()
	if err != nil {
		ctx.Abort(500, GetError)
		if Settings.Debug {
			ctx.WriteString("\r\n")
			debug := fmt.Sprintf("Debug: %s", err)
			ctx.WriteString(debug)
		}
		h.logger.Printf("Dequeue from <%s> Error:%s", val, err)
		return
	}
	resp, _ := json.Marshal(mesg)
	//mesg.(type) is iteface{}
	//resp.(type) is []byte

	ctx.SetHeader("Content-Type", "application/json; charset=UTF-8", true)
	ctx.WriteString(string(resp))

}

func (h *RestQueueHandler) Put(ctx *web.Context, val string) {
	queue := h.Queue(val)
	if !queue.Exists() {
		h.logger.Printf("Queue [%s] didn't existst, will be ceated.", val)
	}
	if mesg, ok := ctx.Params["value"]; ok {
		var i interface{}
		err := json.Unmarshal([]byte(mesg), &i)
		if err != nil {
			ctx.ResponseWriter.WriteHeader(400)
			ctx.WriteString(JsonDecodeError)
			return
		}
		err = queue.Put(i)
		if err != nil {
			ctx.Abort(500, PostError)
			if Settings.Debug {
				ctx.WriteString("\r\n")
				debug := fmt.Sprintf("Debug: %s", err)
				ctx.WriteString(debug)
			}
			h.logger.Printf("Post message into [%s] Error:%s", val, err)
			return
		}
		h.logger.Printf("Put message into queue [%s]", val)

	} else {
		ctx.ResponseWriter.WriteHeader(400)
		ctx.WriteString(LackPostValue)

	}

}

func (h *RestQueueHandler) Clear(ctx *web.Context, val string) {
	queue := h.Queue(val)
	if !queue.Exists() {
		ctx.NotFound(QueueNotFound)
		return
	}
	err := queue.Clear()
	if err != nil {
		ctx.Abort(500, ClearError)
		if Settings.Debug {
			ctx.WriteString("\r\n")
			debug := fmt.Sprintf("Debug: %s", err)
			ctx.WriteString(debug)
		}
		h.logger.Printf("Delete queue [%s] Error:%s", val, err)
		return
	}
	h.logger.Printf("Queue [%s] deleted sucess", val)
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
	web.Get("/q", queueHandler.List)
	web.Get("/q/(.+)", queueHandler.Get)
	web.Post("/q/(.+)", queueHandler.Put)
	web.Delete("/q/(.+)", queueHandler.Clear)
	web.SetLogger(logger)
	web.Run(Settings.HTTPServer.Addr())
}
