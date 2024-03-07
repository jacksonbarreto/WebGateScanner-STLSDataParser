package processor

type IProcessor interface {
	ProcessFile(filePath string) error
}
