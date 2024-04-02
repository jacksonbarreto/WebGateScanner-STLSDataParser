package processor

import (
	"github.com/jacksonbarreto/WebGateScanner-STLSDataParser/config"
	"github.com/jacksonbarreto/WebGateScanner-STLSDataParser/internal/models"
	"github.com/jacksonbarreto/WebGateScanner-STLSDataParser/internal/sslresponseparser"
	kmodels "github.com/jacksonbarreto/WebGateScanner-kafka/models"
	"github.com/jacksonbarreto/WebGateScanner-kafka/producer"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type DefaultFileProcessor struct {
	brokerProducer       producer.IProducer
	fileParser           sslresponseparser.AssessmentSSLResponseParser
	hostExtractor        HostExtractor
	filesInProcess       map[string]bool
	lock                 *sync.Mutex
	readyToProcessSuffix string
}

func NewDefaultFileProcessor(brokerProducer producer.IProducer,
	fileParser sslresponseparser.AssessmentSSLResponseParser,
	hostExtractor HostExtractor, readyToProcessSuffix string, lock *sync.Mutex,
	filesInProcess map[string]bool) *DefaultFileProcessor {
	return &DefaultFileProcessor{
		brokerProducer:       brokerProducer,
		fileParser:           fileParser,
		hostExtractor:        hostExtractor,
		filesInProcess:       filesInProcess,
		lock:                 lock,
		readyToProcessSuffix: readyToProcessSuffix,
	}
}

func (dfp *DefaultFileProcessor) ProcessFileFromChannel(files <-chan string) {
	for filePath := range files {
		log.Println("Starting Processing file: ", filePath)
		dfp.lock.Lock()
		if beingProcessed, exists := dfp.filesInProcess[filePath]; exists && !beingProcessed {
			dfp.filesInProcess[filePath] = true
			log.Println("Processing file (set true): ", filePath)
			dfp.lock.Unlock()
		} else {
			log.Println("File already being processed: ", filePath)
			dfp.lock.Unlock()
			continue
		}

		fileContent, readingFileError := os.ReadFile(filepath.Join(config.App().PathToWatch, filePath))
		if readingFileError != nil {
			dfp.deleteFileFromProcess(filePath)
			log.Println("Error reading file: ", readingFileError)
			// TODO: log error
			continue
		}

		assessmentResult, parsingError := dfp.fileParser.ParseFromJSON(fileContent)
		if parsingError != nil {
			dfp.deleteFileFromProcess(filePath)
			log.Println("Error parsing file: ", parsingError)
			// TODO: log error
			continue
		}

		publishError := dfp.publishToBroker(assessmentResult, filePath)
		if publishError != nil {
			renameError := os.Rename(filePath, filepath.Join(config.App().ErrorParsePath, filepath.Base(filePath)))
			if renameError != nil {
				log.Println("Error moving file to error path: ", renameError)
				// TODO: log error
			}
			dfp.deleteFileFromProcess(filePath)
			log.Println("Error publishing to broker: ", publishError)
			// TODO: log error
			continue
		}

		removeFileError := os.Remove(filepath.Join(config.App().PathToWatch, filePath))
		if removeFileError != nil {
			dfp.deleteFileFromProcess(filePath)
			log.Println("Error removing file: ", removeFileError)
			// TODO: log error
			continue
		}

		removeDoneFileError := os.Remove(filepath.Join(config.App().PathToWatch, filePath+dfp.readyToProcessSuffix))
		if removeDoneFileError != nil {
			dfp.deleteFileFromProcess(filePath)
			log.Println("Error removing done file: ", removeDoneFileError)
			// TODO: log error
			continue
		}

		dfp.deleteFileFromProcess(filePath)
	}
}

func (dfp *DefaultFileProcessor) deleteFileFromProcess(filePath string) {
	dfp.lock.Lock()
	log.Println("Deleting file from process: ", filePath)
	delete(dfp.filesInProcess, filePath)
	log.Println("Deleted file from process: ", filePath)
	dfp.lock.Unlock()
}

func (dfp *DefaultFileProcessor) publishToBroker(assessmentResult models.TestSSLResult, filePath string) error {
	kafkaAssessmentMessage, messageCreationError := kmodels.CreateKafkaEvaluationResponseMessage(
		dfp.hostExtractor.ExtractHostFromFilePath(filePath),
		config.App().Id,
		assessmentResult.StartTime,
		assessmentResult.EndTime,
		assessmentResult)
	if messageCreationError != nil {
		return messageCreationError
	}

	_, _, producerError := dfp.brokerProducer.SendMessage(kafkaAssessmentMessage)
	if producerError != nil {
		return producerError
	}
	return nil
}
