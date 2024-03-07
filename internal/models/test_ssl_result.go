package models

type TestSSLResult struct {
	Endpoints []Endpoint
}

type Endpoint struct {
	IpAddress string
}
