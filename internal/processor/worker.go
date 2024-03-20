package processor

import (
	"github.com/jacksonbarreto/WebGateScanner-STLSDataParser/internal/parser"
	"github.com/jacksonbarreto/WebGateScanner-kafka/producer"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type Worker struct {
	kafkaProducer  producer.IProducer
	errorParsePath string
	lock           *sync.Mutex
	psr            parser.IParser
	processing     map[string]bool
}

func NewWorker(kafkaProducer producer.IProducer, errorParsePath string, lock *sync.Mutex, psr parser.IParser,
	processing map[string]bool) *Worker {
	return &Worker{
		kafkaProducer:  kafkaProducer,
		errorParsePath: errorParsePath,
		lock:           lock,
		psr:            psr,
		processing:     processing,
	}
}

func (w *Worker) Do(files <-chan string) {
	for filePath := range files {
		log.Println("Processing file:", filePath)
		process := NewProcessor(w.kafkaProducer, w.psr)
		if err := process.ProcessFile(filePath); err != nil {
			log.Println("Failed to process file:", err)
			err := os.Rename(filePath, filepath.Join(w.errorParsePath, filepath.Base(filePath)))
			if err != nil {
				log.Println("Failed to move file:", err)
			}
		} else {
			err := os.Remove(filePath)
			if err != nil {
				log.Println("Failed to remove file:", err)
			}
			err = os.Remove(filePath + ".done")
			if err != nil {
				log.Println("Failed to remove file:", err)
			}
		}
		w.lock.Lock()
		delete(w.processing, filePath)
		w.lock.Unlock()
	}
}
