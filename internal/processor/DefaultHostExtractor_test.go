package processor

import (
	"testing"
)

func TestEmptyFilePath(t *testing.T) {
	dhe := DefaultHostExtractor{extension: "json"}
	filePath := ""
	expected := ""

	result := dhe.ExtractHostFromFilePath(filePath)

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestBaseNameWithoutExtension(t *testing.T) {
	dhe := DefaultHostExtractor{extension: ".json"}
	filePath := "/home/user/Downloads/ipb.pt"
	expected := "ipb.pt"

	result := dhe.ExtractHostFromFilePath(filePath)

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestBaseNameWithOnlyDot(t *testing.T) {
	dhe := DefaultHostExtractor{extension: "json"}
	filePath := "/home/user/Downloads/."
	expected := ""

	result := dhe.ExtractHostFromFilePath(filePath)

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestBaseNameWithOnlyTowDot(t *testing.T) {
	dhe := DefaultHostExtractor{extension: "json"}
	filePath := "/home/user/Downloads/.."
	expected := ""

	result := dhe.ExtractHostFromFilePath(filePath)

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestRootFilePath(t *testing.T) {
	dhe := DefaultHostExtractor{extension: "json"}
	filePath := "/"
	expected := ""

	result := dhe.ExtractHostFromFilePath(filePath)

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestEmptyBaseName(t *testing.T) {
	dhe := DefaultHostExtractor{extension: "json"}
	filePath := "/home/user/Downloads/"
	expected := ""

	result := dhe.ExtractHostFromFilePath(filePath)

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}
