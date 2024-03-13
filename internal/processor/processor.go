package processor

import (
	"encoding/json"
	"github.com/jacksonbarreto/WebGateScanner-STLSDataParser/config"
	"github.com/jacksonbarreto/WebGateScanner-STLSDataParser/internal/parser"
	kmodels "github.com/jacksonbarreto/WebGateScanner-kafka/models"
	"github.com/jacksonbarreto/WebGateScanner-kafka/producer"
	stls "github.com/jacksonbarreto/WebGateScanner-stls/models"
	"os"
	"path/filepath"
	"time"
)

type Processor struct {
	KafkaProducer producer.IProducer
	psr           parser.IParser
	// I will need include logger
}

func NewProcessor(producer producer.IProducer, parsing parser.IParser) *Processor {
	return &Processor{
		KafkaProducer: producer,
		psr:           parsing,
	}
}

func (p *Processor) ProcessFile(filePath string) error {
	startTime := time.Now().Unix()
	fileContent, readFileErr := os.ReadFile(filePath)
	if readFileErr != nil {
		return readFileErr
	}

	var response stls.TestSSLResponse
	if unmarshalErr := json.Unmarshal(fileContent, &response); unmarshalErr != nil {
		return unmarshalErr
	}

	result, parseErr := p.psr.ParseJson(response)
	if parseErr != nil {
		return parseErr
	}

	endTime := time.Now().Unix()
	// TODO: change the start time and end time to get from the result
	// TODO: change the institutionID to get from the result
	kafkaMessage, errMessage := kmodels.CreateKafkaEvaluationResponseMessage("url", config.App().Id,
		startTime, endTime, result)
	if errMessage != nil {
		return errMessage

	}

	_, _, producerErr := p.KafkaProducer.SendMessage(kafkaMessage)
	if producerErr != nil {
		renameErr := os.Rename(filePath, filepath.Join(config.App().ErrorParsePath, filepath.Base(filePath)))
		if renameErr != nil {
			return renameErr
		}
		return producerErr
	}

	removeErr := os.Remove(filePath)
	if removeErr != nil {
		return removeErr
	}

	return nil
}
