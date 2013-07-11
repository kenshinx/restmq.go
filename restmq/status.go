package restmq

const (
	QueueNotFound      = "Queue not Found"
	EmptyQueue         = "Queue is Empty"
	GetError           = "Get message from queue failed"
	LackPostValue      = "PoInvalid json formatting datast params must be consist of 'value: '"
	JsonDecodeError    = "Invalid json formatting data"
	PostError          = "Post message into queue failed"
	ClearError         = "Delete Queue failed"
	WebSocketConnError = "must be a websocket handshake"
	ConsumeError       = "Consume from queue failed"
)
