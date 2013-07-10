# RestMQ

RestMQ is a message queue based on redis.
Multiple-type data manipulation ways supported: 

* HTTP GET/POST/PUT/DELETE
* WebSocket
* Comet

This repo is rewritten by Go. and there is a completely implements with python:
  [https://github.com/gleicon/restmq](https://github.com/gleicon/restmq)


### Restful API

* GET Message 

```
curl -i  http://localhost:8000/q/redisq:kenshix
```

* PUT Message 

```
curl -i -X POST http://localhost:8000/q/kenshin -d 'value={"x":1,"y":2}'  
curl -i -X POST http://localhost:8000/q/kenshin -d 'value="xx"'
curl -i -X POST http://localhost:8000/q/kenshin -d 'value=1'
```

1. The post params must conatin value field.
2. The message should be json encode

