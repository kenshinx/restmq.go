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
	"time"
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
			h.logger.Printf("Put message into [%s] Error:%s", val, err)
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

const (
	// Time allowed to write a message to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next message from the client.
	readWait = 60 * time.Second

	// Send pings to client with this period. Must be less than readWait.
	pingPeriod = (readWait * 9) / 10

	//Consume wait should bigger than ping period
	consumeWait = pingPeriod * 2

	// Maximum message size allowed from client.
	maxMessageSize = 512
)

type WSQueueHandler struct {
	redis  *redis.Client
	logger *log.Logger
}

func (wsh *WSQueueHandler) Consumer(ctx *web.Context, val string) {

	ws, err := wsh.handshake(ctx.Request, ctx.ResponseWriter)
	if err != nil {
		writeError(ctx, 400, WebSocketConnError)
	}

	queue := redisq.NewRedisQueue(wsh.redis, val)
	if !queue.Exists() {
		ws.WriteMessage(websocket.OpText, []byte(QueueNotFound))
		ws.WriteMessage(websocket.OpClose, []byte{})
		return
	}

	wsh.logger.Printf("Get websocket connection from %s", ws.RemoteAddr())
	wsh.logger.Printf("Begin subscribe queue [%s]", val)

	c := WebSocketConn{ws, queue, wsh.logger}
	go c.writePump()
	// c.readPump()

}

func (wsh *WSQueueHandler) handshake(r *http.Request, w http.ResponseWriter) (ws *websocket.Conn, err error) {
	ws, err = websocket.Upgrade(w, r.Header, nil, 1024, 1024)
	return
}

type WebSocketConn struct {
	ws     *websocket.Conn
	rq     *redisq.RedisQueue
	logger *log.Logger
}

func (c *WebSocketConn) write(opCode int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(opCode, payload)
}

func (c *WebSocketConn) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	var mesgs = make(chan interface{})
	c.rq.Consume(true, uint(consumeWait.Seconds()), mesgs)
	for {
		select {
		case v, ok := <-mesgs:
			if !ok {
				c.write(websocket.OpText, []byte(ConsumeError))
				c.write(websocket.OpClose, []byte{})
				c.logger.Printf("Consumer from %s failed: %s", c.rq)
				return
			}
			mesg, _ := json.Marshal(v)
			if err := c.write(websocket.OpText, mesg); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.OpPing, []byte{}); err != nil {
				return
			}

		}
	}

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
