package models

type MapPageDataModel struct {
	Services []ServiceModel `json:"services"`
	Tags     []string       `json:"tags"`
}
