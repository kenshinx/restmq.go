package restmq

const (
	QueueNotFound      = "Queue not Found"
	EmptyQueue         = "Queue is Empty"
	GetError           = "Get message from queue failed"
	LackPostValue      = "Post params must be consist of 'value: '"
	JsonDecodeError    = "Invalid json formatting data"
	PostError          = "Post message into queue failed"
	ClearError         = "Delete Queue failed"
	WebSocketConnError = "Must be a websocket handshake"
	ConsumeError       = "Consume from queue failed"
)
