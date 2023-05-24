package rabbit

import (
	"bytes"
	"encoding/json"

	"github.com/ManuelP84/calendar_notification/domain/task/events"
)

func Deserialize(b []byte) (events.TaskEvent, error) {
	var taskEvent events.TaskEvent
	buf := bytes.NewBuffer(b)
	decoder := json.NewDecoder(buf)
	err := decoder.Decode(&taskEvent)
	return taskEvent, err
}
