# RestMQ

RestMQ is a message queue based on redis.
Multiple-type data manipulation ways supported: 

* HTTP GET/POST/PUT/DELETE
* WebSocket
* Comet

This repo is rewritten by Go. and there is a completely implements with python:
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




