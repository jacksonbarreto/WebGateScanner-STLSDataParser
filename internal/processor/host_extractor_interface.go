package processor

type HostExtractor interface {
	ExtractHostFromFilePath(filePath string) string
}
