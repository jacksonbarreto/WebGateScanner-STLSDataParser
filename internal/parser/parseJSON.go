package parser

import (
	"github.com/jacksonbarreto/WebGateScanner-STLSDataParser/internal/models"
	stls "github.com/jacksonbarreto/WebGateScanner-stls/models"
)

type Parser struct {
	// I will need include logger
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParseJson(response stls.TestSSLResponse) (models.TestSSLResult, error) {
	var result = models.TestSSLResult{}

	result.Endpoints = make([]models.Endpoint, len(response.ScanResult))
	for i, endpoint := range response.ScanResult {
		result.Endpoints[i] = models.Endpoint{
			IpAddress: endpoint.IP,
		}
	}

	return result, nil
}
