# RestMQ

RestMQ is a message queue based on redis.
Support Multiple-type data manipulation ways. 

* HTTP GET/POST/PUT/DELETE
* WebSocket
* ~~Comet~~  not yet

This repo is rewritten by Go. and there is a more completely implements with python:
  [https://github.com/gleicon/restmq](https://github.com/gleicon/restmq)



### Install & Run

```
$ go get github.com/kenshinx/restmq.go
$ cd $GOPATH/src/github.com/kenshinx/restmq.go
$ go run main.go -c restmq.conf
```

### Configuration

All Configuration in the restmq.conf file.  This config file support Toml rule.  

More about Toml :[https://github.com/mojombo/toml](https://github.com/mojombo/toml)



### Restful API

* List Queues 

```
curl -i -X GET http://localhost:8000/q
```

* Get Message 

```
curl -i  http://localhost:8000/q/kenshin
```

* Put Message 

```
curl -i -X POST http://localhost:8000/q/kenshin -d 'value={"x":1,"y":2}'  
curl -i -X POST http://localhost:8000/q/kenshin -d 'value="xx"'
curl -i -X POST http://localhost:8000/q/kenshin -d 'value=1'
```

The post message should be json encoded

* Delete Queue

```
curl -i -X DELETE http://localhost:8000/q/kenshin
```

### WebSocket Consumer

*The Consumer Url*

```
http://localhost:8000/ws/kenshin
```

 The client demo can reference `websocket.html`


## Dependence

* [github.com/kenshinx/redisq](https://github.com/kenshinx/redisq)
* [github.com/hoisie/redis](https://github.com/hoisie/redis)
* [github.com/hoisie/web](https://github.com/hoisie/web)
* [github.com/garyburd/go-websocket/websocket](https://github.com/garyburd/go-websocket/websocket)
* [github.com/BurntSushi/toml](https://github.com/BurntSushi/toml)


##TODO
* Control protocol
* Comet consumer
* Web dashboard















