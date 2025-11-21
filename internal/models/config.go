package models

type Config struct {
	TemplateName     string            `json:"templateName" bson:"templateName"`
	Variables        map[string]string `json:"variables" bson:"variables"`
	PackageDirectory string            `json:"packageDirectory" bson:"packageDirectory"`
}
