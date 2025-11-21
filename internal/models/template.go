package models

type Template struct {
	Name        string              `json:"name" bson:"name"`
	Description string              `json:"description" bson:"description"`
	Variables   map[string]string   `json:"variables" bson:"variables"`
	Detect      map[string]string   `json:"detect" bson:"detect"`
	Steps       []map[string]string `json:"steps" bson:"steps"`
}
