package models

// import (
// 	"reflect"
// 	"testing"
// )
//
// func TestContructLocationFromLatlong(t *testing.T) {
// 	expected := "POINT(14.123456 121.123456)"
//
// 	model := &GeoSpatialModel{
// 		Latitude:  14.123456,
// 		Longitude: 121.123456,
// 	}
//
// 	ConstructLocationFromLatLong(model)
//
// 	if model.Location != expected {
// 		t.Error("Expected: ", expected)
// 		t.Error("Got: ", model.Location)
// 		t.FailNow()
// 	}
// }
//
// func TestExtractLatLongFromLocation(t *testing.T) {
// 	expected := &GeoSpatialModel{
// 		Latitude:  14.123456,
// 		Longitude: 121.123456,
// 		Location:  "POINT(14.123456 121.123456)",
// 	}
//
// 	testModel := &GeoSpatialModel{
// 		Location: "POINT(14.123456 121.123456)",
// 	}
//
// 	err := ExtractLatLongFromLocation(testModel)
// 	if err != nil {
// 		t.Error("Expected: no error")
// 		t.Error("Got: Error: ", err.Error())
// 	}
//
// 	if reflect.DeepEqual(testModel, expected) == false {
// 		t.Error("Expected: ", expected)
// 		t.Error("Got: ", testModel)
// 		t.FailNow()
// 	}
// }
