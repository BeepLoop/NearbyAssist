package filehandler

import (
	"testing"
	"time"
)

func TestCreateServicePhotoTimestamp(t *testing.T) {
	expected := time.Now().Format("2006-01-02_15:04:05")
	result := NewServicePhoto(1, 1, nil)

	if result.Timestamp != expected {
		t.Fatalf("Expected %s but got %s", expected, result.Timestamp)
	}
}

func TestCreateServicePhoto(t *testing.T) {
	result := NewServicePhoto(1, 1, nil)

	if result.VendorId != 1 {
		t.Fatalf("Expected 1 but got %d", result.VendorId)
	}

	if result.ServiceId != 1 {
		t.Fatalf("Expected 1 but got %d", result.ServiceId)
	}

	if result.File != nil {
		t.Fatalf("Expected nil but got %v", result.File)
	}
}
