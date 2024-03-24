package processor

type FileProcessor interface {
	ProcessFileFromChannel(files <-chan string)
}
