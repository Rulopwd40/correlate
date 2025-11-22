package pipeline

type EventType string

const (
	EventTaskStart    EventType = "task_start"
	EventTaskProgress EventType = "task_progress"
	EventTaskFinish   EventType = "task_finish"
	EventPipelineDone EventType = "pipeline_done"
	EventError        EventType = "error"
)

type Event struct {
	Type       EventType
	TaskName   string
	Message    string
	Percentage int
	Err        error
}
