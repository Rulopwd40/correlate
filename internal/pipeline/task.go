package pipeline

type Task struct {
	Cmd     string `json:"cmd" bson:"cmd"`
	Name    string `json:"name" bson:"name"`
	Workdir string `json:"workdir" bson:"workdir"`
}
