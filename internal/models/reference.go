package models

type Reference struct {
	Identifier          string   `json:"identifier" bson:"identifier"`
	ManifestDirectories []string `json:"directories" bson:"directories"`
}

type References struct {
	References []Reference `json:"references" bson:"references"`
}
