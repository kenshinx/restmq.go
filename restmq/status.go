package restmq

import (
	"fmt"
)

var Status = QueueStatus{}

type QueueStatus struct {
}

func (s QueueStatus) QueueNotFound(queue string) (mesg string) {
	mesg = fmt.Sprintf("Queue [%s] not found\n", queue)
	return
}

func (s QueueStatus) EmptyQueue(queue string) (mesg string) {
	mesg = fmt.Sprintf("Queue [%s] is Empty\n", queue)
	return
}

func (s QueueStatus) GetError(queue string) (mesg string) {
	mesg = fmt.Sprintf("Get message from queue [%s] failed\n", queue)
	return
}

func (s QueueStatus) LackPostValue() (mesg string) {
	return "Post params lack of value filed\n"
}

func (s QueueStatus) JsonDecodeError() (mesg string) {
	return "Invalid json formatting data\n"
}

func (s QueueStatus) PostError() (mesg string) {
	return "Post message into queue failed\n"
}
