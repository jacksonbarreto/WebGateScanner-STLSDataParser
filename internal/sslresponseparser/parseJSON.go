package sslresponseparser

import (
	"encoding/json"
	"fmt"
	"github.com/jacksonbarreto/WebGateScanner-STLSDataParser/internal/models"
	stls "github.com/jacksonbarreto/WebGateScanner-stls/models"
	"strconv"
)

type DefaultAssessmentSSLResponseParser struct {
	// I will need include logger
}

func NewDefaultAssessmentSSLResponseParser() *DefaultAssessmentSSLResponseParser {
	return &DefaultAssessmentSSLResponseParser{}
}

func (p *DefaultAssessmentSSLResponseParser) ParseFromJSON(fileContent []byte) (models.TestSSLResult, error) {
	var response stls.TestSSLResponse
	if unmarshalError := json.Unmarshal(fileContent, &response); unmarshalError != nil {
		return models.TestSSLResult{}, unmarshalError
	}

	startTimeInt, err := p.stringTimeToUnixTime(response.StartTime)
	if err != nil {
		return models.TestSSLResult{}, err
	}

	var result = models.TestSSLResult{
		StartTime:     startTimeInt,
		EndTime:       startTimeInt + int64(response.ScanTime),
		Endpoints:     make([]models.Endpoint, len(response.ScanResult)),
		RawAssessment: string(fileContent),
	}

	for i, endpoint := range response.ScanResult {
		result.Endpoints[i] = models.Endpoint{
			IpAddress:  endpoint.IP,
			TargetHost: endpoint.TargetHost,
			Protocols:  processSSLProtocols(endpoint.Protocols),
		}
	}

	return result, nil
}

func (p *DefaultAssessmentSSLResponseParser) stringTimeToUnixTime(time string) (int64, error) {
	startTimeInt, err := strconv.ParseInt(time, 10, 64)
	if err != nil {
		fmt.Println("Erro ao converter StartTime:", err)
		return 0, err
	}
	return startTimeInt, nil
}

func processSSLProtocols(protocols []stls.Protocols) []models.Protocol {
	result := make([]models.Protocol, len(protocols))
	for _, protocol := range protocols {
		if protocol.Finding == "offered" || protocol.Finding == "offered with final" {
			result = append(result, models.Protocol{
				Name: protocol.ID,
			})
		}
	}
	return result
}
