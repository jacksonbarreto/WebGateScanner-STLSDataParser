package parser

import (
	"github.com/jacksonbarreto/WebGateScanner-STLSDataParser/internal/models"
	stls "github.com/jacksonbarreto/WebGateScanner-stls/models"
)

type IParser interface {
	ParseJson(response stls.TestSSLResponse) (models.TestSSLResult, error)
}
