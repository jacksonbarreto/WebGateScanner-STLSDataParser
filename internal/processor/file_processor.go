package processor

import (
	"github.com/jacksonbarreto/WebGateScanner-STLSDataParser/config"
	"github.com/jacksonbarreto/WebGateScanner-STLSDataParser/internal/models"
	"github.com/jacksonbarreto/WebGateScanner-STLSDataParser/internal/sslresponseparser"
	kmodels "github.com/jacksonbarreto/WebGateScanner-kafka/models"
	"github.com/jacksonbarreto/WebGateScanner-kafka/producer"
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
		dfp.lock.Lock()
		if beingProcessed, exists := dfp.filesInProcess[filePath]; exists && !beingProcessed {
			dfp.filesInProcess[filePath] = true
			dfp.lock.Unlock()
		} else {
			dfp.lock.Unlock()
			continue
		}

		fileContent, readingFileError := os.ReadFile(filePath)
		if readingFileError != nil {
			dfp.deleteFileFromProcess(filePath)
			// TODO: log error
			continue
		}

		assessmentResult, parsingError := dfp.fileParser.ParseFromJSON(fileContent)
		if parsingError != nil {
			dfp.deleteFileFromProcess(filePath)
			// TODO: log error
			continue
		}

		publishError := dfp.publishToBroker(assessmentResult, filePath)
		if publishError != nil {
			renameError := os.Rename(filePath, filepath.Join(config.App().ErrorParsePath, filepath.Base(filePath)))
			if renameError != nil {
				// TODO: log error
			}
			dfp.deleteFileFromProcess(filePath)
			// TODO: log error
			continue
		}

		removeFileError := os.Remove(filePath)
		if removeFileError != nil {
			dfp.deleteFileFromProcess(filePath)
			// TODO: log error
			continue
		}

		removeDoneFileError := os.Remove(filePath + dfp.readyToProcessSuffix)
		if removeDoneFileError != nil {
			dfp.deleteFileFromProcess(filePath)
			// TODO: log error
			continue
		}

		dfp.deleteFileFromProcess(filePath)
	}
}

func (dfp *DefaultFileProcessor) deleteFileFromProcess(filePath string) {
	dfp.lock.Lock()
	delete(dfp.filesInProcess, filePath)
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
