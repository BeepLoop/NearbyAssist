package filehandler

import (
	"testing"
	"time"
)

func TestCreateApplicationProofTimestamp(t *testing.T) {
	expected := time.Now().Format("2006-01-02_15:04:05")
	result := NewApplicationProof(1, 1, nil)

	if result.Timestamp != expected {
		t.Fatalf("Expected %s but got %s", expected, result.Timestamp)
	}
}

func TestCreateApplicationProof(t *testing.T) {
	result := NewApplicationProof(1, 1, nil)

	if result.ApplicationId != 1 {
		t.Fatalf("Expected 1 but got %d", result.ApplicationId)
	}

	if result.ApplicantId != 1 {
		t.Fatalf("Expected 1 but got %d", result.ApplicantId)
	}

	if result.File != nil {
		t.Fatalf("Expected nil but got %v", result.File)
	}
}
