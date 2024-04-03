package sslresponseparser

import (
	"encoding/json"
	"github.com/jacksonbarreto/WebGateScanner-STLSDataParser/internal/models"
	stls "github.com/jacksonbarreto/WebGateScanner-stls/models"
	"log"
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
		StartTime:      startTimeInt,
		EndTime:        startTimeInt + int64(response.ScanTime),
		EngineVersion:  response.Version,
		OpenSSLVersion: response.Openssl,
		Endpoints:      make([]models.Endpoint, len(response.ScanResult)),
		RawAssessment:  string(fileContent),
	}

	for i, endpoint := range response.ScanResult {
		finalScore, errScore := getFinalScore(endpoint.Rating)
		if errScore != nil || finalScore == -1 {
			log.Printf("Error getting final score from endpoint(%s): %v\n", endpoint.IP, errScore)
		}
		result.Endpoints[i] = models.Endpoint{
			IpAddress:          endpoint.IP,
			TargetHost:         endpoint.TargetHost,
			FinalScore:         finalScore,
			OverallGrade:       getOverallGrade(endpoint.Rating),
			Protocols:          getSSLProtocols(endpoint.Protocols),
			DigitalCertificate: getDigitalCertificate(endpoint.ServerDefaults),
		}
	}

	return result, nil
}

func (p *DefaultAssessmentSSLResponseParser) stringTimeToUnixTime(time string) (int64, error) {
	startTimeInt, err := strconv.ParseInt(time, 10, 64)
	if err != nil {
		log.Println("Erro ao converter StartTime:", err)
		return 0, err
	}
	return startTimeInt, nil
}

func getSSLProtocols(protocols []stls.Protocols) []models.Protocol {
	var result []models.Protocol
	for _, protocol := range protocols {
		if protocol.Finding != "not offered" {
			result = append(result, models.Protocol{
				Name: protocol.ID,
			})
		}
	}
	return result
}

func getFinalScore(Rating []stls.Rating) (int, error) {
	for _, rating := range Rating {
		if rating.ID == "final_score" {
			score, err := strconv.Atoi(rating.Finding)
			if err != nil {
				return 0, err
			}
			return score, nil
		}
	}
	return -1, nil
}

func getOverallGrade(Rating []stls.Rating) string {
	for _, rating := range Rating {
		if rating.ID == "overall_grade" {
			return rating.Finding
		}
	}
	return ""
}

func getDigitalCertificate(ServerDefaults []stls.ServerDefaults) models.DigitalCertificate {
	var result models.DigitalCertificate
	for _, serverDefault := range ServerDefaults {
		if serverDefault.ID == "cert_signatureAlgorithm" {
			result.SignatureAlgorithm = serverDefault.Finding
		}
		if serverDefault.ID == "cert_keySize" {
			result.KeySize = serverDefault.Finding
		}
		if serverDefault.ID == "cert_notBefore" {
			result.NotBefore = serverDefault.Finding
		}
		if serverDefault.ID == "cert_notAfter" {
			result.NotAfter = serverDefault.Finding
		}
		if serverDefault.ID == "cert_caIssuers" {
			result.CaIssuers = serverDefault.Finding
		}
		if serverDefault.ID == "cert_chain_of_trust" {
			result.ChainOfTrustPassed = serverDefault.Finding == "passed."
		}
		if serverDefault.ID == "DNS_CAArecord" {
			result.DnsCAA = serverDefault.Finding
		}
		if serverDefault.ID == "cert_commonName" {
			result.CommonName = serverDefault.Finding
		}
		if serverDefault.ID == "cert_subjectAltName" {
			result.SubjectAltName = serverDefault.Finding
		}
		if serverDefault.ID == "cert_commonName_wo_SNI" {
			result.CommonNameWoSNI = serverDefault.Finding
		}
	}
	return result
}
