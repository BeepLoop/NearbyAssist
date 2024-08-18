package models

import "testing"

func TestNewServicePhotoModel(t *testing.T) {
	servicePhotoId := "1"
	testVendorId := "1"
	testServiceId := "1"
	testFilename := "filename1.jpeg"
	expectedUrl := "/resource/service/filename1.jpeg"

	model := NewServicePhotoModel(servicePhotoId, testVendorId, testServiceId, testFilename)

	if model.VendorId != testVendorId {
		t.Fatalf("Expected: %s, Got: %s", testVendorId, model.VendorId)
	}

	if model.ServiceId != testServiceId {
		t.Fatalf("Expected: %s, Got: %s", testServiceId, model.ServiceId)
	}

	if model.Url != expectedUrl {
		t.Fatalf("Expected: %s, Got: %s", expectedUrl, model.Url)
	}
}
