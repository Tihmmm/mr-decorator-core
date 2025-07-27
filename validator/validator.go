package validator

import (
	"github.com/Tihmmm/mr-decorator-core/models"
	"github.com/Tihmmm/mr-decorator-core/parser"
	jsonvalidator "github.com/go-playground/validator/v10"
	"log"
	"slices"
)

type Validator interface {
	IsValidAll(reqBody *models.MRRequest) bool
	IsValidStruct(reqBody *models.MRRequest) bool
	isValidFormat(fileName string) bool
}

type RequestValidator struct {
	jsonValidator     *jsonvalidator.Validate
	registeredFormats []string
}

func NewValidator() Validator {
	return &RequestValidator{
		jsonValidator:     jsonvalidator.New(),
		registeredFormats: parser.List(),
	}
}

func (v *RequestValidator) IsValidAll(reqBody *models.MRRequest) bool {
	return v.IsValidStruct(reqBody) && v.isValidFormat(reqBody.ArtifactFileName)
}

func (v *RequestValidator) IsValidStruct(reqBody *models.MRRequest) bool {
	err := v.jsonValidator.Struct(reqBody)
	if err == nil {
		return true
	}

	log.Printf("Error validating struct: %s\n", err)
	return false
}

func (v *RequestValidator) isValidFormat(format string) bool {
	if slices.Contains(v.registeredFormats, format) {
		return true
	}

	log.Printf("Invalid artifact format: %s\n", format)
	return false
}
