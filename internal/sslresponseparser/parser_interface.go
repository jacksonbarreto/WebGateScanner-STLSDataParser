package sslresponseparser

import (
	"github.com/jacksonbarreto/WebGateScanner-STLSDataParser/internal/models"
)

type AssessmentSSLResponseParser interface {
	ParseFromJSON(fileContent []byte) (models.TestSSLResult, error)
}
