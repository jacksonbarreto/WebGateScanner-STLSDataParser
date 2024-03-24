package models

type TestSSLResult struct {
	StartTime     int64
	EndTime       int64
	Endpoints     []Endpoint
	RawAssessment string
}

type Endpoint struct {
	IpAddress  string
	TargetHost string
	Protocols  []Protocol
}

type Protocol struct {
	Name string
}
