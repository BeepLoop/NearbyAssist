package models

type MapPageDataModel struct {
	Markers []GeoSpatialModel `json:"markers"`
	Tags    []string          `json:"tags"`
}
