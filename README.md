# RestMQ

RestMQ is a message queue based on redis.
Multiple-type data manipulation ways supported: 

* HTTP GET/POST/PUT/DELETE
* WebSocket
* ~~Comet~~  not yet

This repo is rewritten by Go. and there is a more completely implements with python:
  [https://github.com/gleicon/restmq](https://github.com/gleicon/restmq)


### Restful API

* List Queues 

```
curl -i -X GET http://localhost:8000/q
```

* Get Message 

```
curl -i  http://localhost:8000/q/redisq:kenshix
```

* Put Message 

```
curl -i -X POST http://localhost:8000/q/kenshin -d 'value={"x":1,"y":2}'  
curl -i -X POST http://localhost:8000/q/kenshin -d 'value="xx"'
curl -i -X POST http://localhost:8000/q/kenshin -d 'value=1'
```

The post message should be json encode

* Delete Queue

```
curl -i -X DELETE http://localhost:8000/q/kenshin
```

### WebSocket Consumer

*The Consumer Url*

```
http://localhost:8000/ws/kenshin
```


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















