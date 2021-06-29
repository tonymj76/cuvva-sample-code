package models

import "github.com/xeipuuv/gojsonschema"

type (
	CreateRequest struct {
		NumberOfProduct int64  `json:"number_of_product"`
		Email           string `json:"email"`
		BusinessName    string `json:"business_name"`
	}

	CreateResponse struct {
		ID              uint64 `json:"id"`
		NumberOfProduct int64  `json:"number_of_product"`
		Email           string `json:"email"`
		BusinessName    string `json:"business_name"`
		CreatedAt       string `json:"Created_at"`
	}
)

var CreateRequestSchema = gojsonschema.NewStringLoader(`{
	"$schema": "http://json-schema.org/schema#",
	"type": "object",
	"required": [ "number_of_product", "email", "business_name" ],
	"properties": {
		"email": {
			"type": "string", "format": "email"
		},
		"business_name": {
			"type": "string"
		},
		"number_of_product": {
			"type": "integer"
		}
	}
}`)
