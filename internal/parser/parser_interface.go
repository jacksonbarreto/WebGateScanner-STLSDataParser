package parser

import stls "github.com/jacksonbarreto/WebGateScanner-stls/models"

type IParser interface {
	ParseJson(response stls.TestSSLResponse) (string, error)
}
