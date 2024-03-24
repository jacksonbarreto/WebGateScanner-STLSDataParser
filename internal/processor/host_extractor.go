package processor

import (
	"path/filepath"
	"strings"
)

type DefaultHostExtractor struct {
	extension string
}

func NewDefaultHostExtractor(extension string) *DefaultHostExtractor {
	return &DefaultHostExtractor{
		extension: extension,
	}
}

func (dhe DefaultHostExtractor) ExtractHostFromFilePath(filePath string) string {
	if filePath == "" || strings.HasSuffix(filePath, "/") {
		return ""
	}
	baseName := filepath.Base(filePath)

	if baseName == "." ||
		baseName == ".." ||
		baseName == "/" {
		return ""
	}

	if !strings.HasSuffix(baseName, dhe.extension) {
		return baseName
	}

	return strings.TrimSuffix(baseName, dhe.extension)
}
