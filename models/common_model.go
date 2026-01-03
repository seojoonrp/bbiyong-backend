// models/common_model.go

package models

// Location GeoJSON 구조체
type Location struct {
	Type        string    `bson:"type" json:"type"`
	Coordinates []float64 `bson:"coordinates" json:"coordinates"`
}
