package models

type TestSSLResult struct {
	StartTime      int64
	EndTime        int64
	EngineVersion  string
	OpenSSLVersion string
	Endpoints      []Endpoint
	RawAssessment  string
}

type Endpoint struct {
	IpAddress          string
	TargetHost         string
	FinalScore         int
	OverallGrade       string
	Protocols          []Protocol
	DigitalCertificate DigitalCertificate
}

type Protocol struct {
	Name string
}

type DigitalCertificate struct {
	KeySize            string
	SignatureAlgorithm string
	CaIssuers          string
	ChainOfTrustPassed bool
	NotBefore          string
	NotAfter           string
	DnsCAA             string
	CommonName         string
	CommonNameWoSNI    string
	SubjectAltName     string
}
