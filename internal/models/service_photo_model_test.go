package models

import "testing"

func TestNewServicePhotoModel(t *testing.T) {
	testVendorId := 1
	testServiceId := 1
	testFilename := "filename1.jpeg"
	expectedUrl := "/resource/service/filename1.jpeg"

	model := NewServicePhotoModel(testVendorId, testServiceId, testFilename)

	if model.VendorId != testVendorId {
		t.Fatalf("Expected: %d, Got: %d", testVendorId, model.VendorId)
	}

	if model.ServiceId != testServiceId {
		t.Fatalf("Expected: %d, Got: %d", testServiceId, model.ServiceId)
	}

	if model.Url != expectedUrl {
		t.Fatalf("Expected: %s, Got: %s", expectedUrl, model.Url)
	}
}
