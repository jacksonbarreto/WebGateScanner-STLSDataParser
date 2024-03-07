package processor

import (
	"encoding/json"
	"github.com/jacksonbarreto/WebGateScanner-STLSDataParser/config"
	"github.com/jacksonbarreto/WebGateScanner-STLSDataParser/internal/parser"
	"github.com/jacksonbarreto/WebGateScanner-kafka/producer"
	stls "github.com/jacksonbarreto/WebGateScanner-stls/models"
	"os"
	"path/filepath"
)

type Processor struct {
	KafkaProducer *producer.IProducer
	Parser        *parser.IParser
	// I will need include logger
}

func New(producer *producer.IProducer, parser *parser.IParser) *Processor {
	return &Processor{
		KafkaProducer: producer,
		Parser:        parser,
	}
}

func (p *Processor) ProcessFile(filePath string) error {
	fileContent, readFileErr := os.ReadFile(filePath)
	if readFileErr != nil {
		return readFileErr
	}

	var response stls.TestSSLResponse
	if unmarshalErr := json.Unmarshal(fileContent, &response); unmarshalErr != nil {
		return unmarshalErr
	}

	result, parseErr := p.Parser.ParseJson(response)
	if parseErr != nil {
		return parseErr
	}

	jsonResponse, marshalErr := json.Marshal(result)
	if marshalErr != nil {
		return marshalErr
	}

	_, _, producerErr := p.KafkaProducer.SendMessage(string(jsonResponse))
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
