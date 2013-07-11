package restmq

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/go-websocket/websocket"
	"github.com/hoisie/redis"
	"github.com/hoisie/web"
	"github.com/kenshinx/redisq"
	"log"
	"net/http"
	"os"
	// "time"
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

func (h *RestQueueHandler) List(ctx *web.Context) {

	var keys []string
	keys, _ = h.redis.Keys("*")
	resp, _ := json.Marshal(keys)
	ctx.SetHeader("Content-Type", "application/json; charset=UTF-8", true)
	ctx.WriteString(string(resp))
}

func (h *RestQueueHandler) Get(ctx *web.Context, val string) {
	queue := redisq.NewRedisQueue(h.redis, val)
	if !queue.Exists() {
		writeError(ctx, 404, QueueNotFound)
		return
	}
	if queue.Empty() {
		writeError(ctx, 400, EmptyQueue)
		return
	}
	mesg, err := queue.GetNoWait()
	if err != nil {
		writeError(ctx, 500, GetError)
		if Settings.Debug {
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
	queue := redisq.NewRedisQueue(h.redis, val)
	if !queue.Exists() {
		h.logger.Printf("Queue [%s] didn't existst, will be ceated.", val)
	}
	if mesg, ok := ctx.Params["value"]; ok {
		var i interface{}
		err := json.Unmarshal([]byte(mesg), &i)
		if err != nil {
			writeError(ctx, 400, JsonDecodeError)
			return
		}
		err = queue.Put(i)
		if err != nil {
			writeError(ctx, 500, PostError)
			if Settings.Debug {
				debug := fmt.Sprintf("Debug: %s", err)
				ctx.WriteString(debug)
			}
			h.logger.Printf("Post message into [%s] Error:%s", val, err)
			return
		}
		h.logger.Printf("Put message into queue [%s]", val)

	} else {
		writeError(ctx, 400, LackPostValue)

	}

}

func (h *RestQueueHandler) Clear(ctx *web.Context, val string) {
	queue := redisq.NewRedisQueue(h.redis, val)
	if !queue.Exists() {
		writeError(ctx, 404, QueueNotFound)
		return
	}
	err := queue.Clear()
	if err != nil {
		writeError(ctx, 500, ClearError)
		if Settings.Debug {
			debug := fmt.Sprintf("Debug: %s", err)
			ctx.WriteString(debug)
		}
		h.logger.Printf("Delete queue [%s] Error:%s", val, err)
		return
	}
	h.logger.Printf("Queue [%s] deleted sucess", val)
}

type WSQueueHandler struct {
	redis  *redis.Client
	logger *log.Logger
}

func (wsh *WSQueueHandler) Consumer(ctx *web.Context, val string) {

	ws, err := wsh.handshake(ctx.Request, ctx.ResponseWriter)
	if err != nil {
		writeError(ctx, 400, WebSocketConnError)
	}
	c := WebSocketConn{ws}
	queue := redisq.NewRedisQueue(wsh.redis, val)
	if !queue.Exists() {
		c.write(websocket.OpText, []byte(QueueNotFound))
		c.write(websocket.OpClose, []byte{})
		return
	}
	wsh.logger.Printf("Get websocket connection from %s", ws.RemoteAddr())
	// wsh.logger.Printf("", ...)
}

func (wsh *WSQueueHandler) handshake(r *http.Request, w http.ResponseWriter) (ws *websocket.Conn, err error) {
	ws, err = websocket.Upgrade(w, r.Header, nil, 1024, 1024)
	return
}

type WebSocketConn struct {
	ws *websocket.Conn
}

func (c *WebSocketConn) write(opCode int, payload []byte) error {
	// c.ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
	return c.ws.WriteMessage(opCode, payload)
}

func writeError(ctx *web.Context, statusCode int, errorMesg string) {
	ctx.ResponseWriter.WriteHeader(statusCode)
	ctx.WriteString(errorMesg)
	ctx.WriteString("\r\n")
}

func initLogger(log_file string) (logger *log.Logger) {
	if log_file != "" {
		f, err := os.Create(log_file)
		if err != nil {
			os.Exit(1)
		}
		logger = log.New(f, "[restmq]", log.Ldate|log.Ltime)
	} else {
		logger = log.New(os.Stdout, "[restmq]", log.Ldate|log.Ltime)
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
		restHandler  = &RestQueueHandler{redis, logger}
		wsHandler    = &WSQueueHandler{redis, logger}
	)

	web.Get("/", indexHandler.Get)
	web.Get("/q", restHandler.List)
	web.Get("/q/(.+)", restHandler.Get)
	web.Post("/q/(.+)", restHandler.Put)
	web.Delete("/q/(.+)", restHandler.Clear)
	web.Get("/ws/(.+)", wsHandler.Consumer)
	web.SetLogger(logger)
	web.Run(Settings.HTTPServer.Addr())
}
