package restmq

import (
	"fmt"
)

var Status = QueueStatus{}

type QueueStatus struct {
}

func (s QueueStatus) QueueNotFound(queue string) (mesg string) {
	mesg = fmt.Sprintf("Queue <%s> not found", queue)
	return
}

func (s QueueStatus) BadRequest(queue string) (mesg string) {
	mesg = fmt.Sprintf("Request queue <%s> failed", queue)
	return
}
