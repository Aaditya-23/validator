package validator

import (
	"fmt"
)

type field interface {
	_parse(*[]Error) bool
}

type Error struct {
	Field   string
	Message string
	Code    string
}

type RefinementData struct {
	Field string
	Code  string
}

func requiredFieldErr(fieldName, required_err string) Error {
	var err Error
	err.Field = fieldName
	err.Code = CodeRequired

	if required_err != "" {
		err.Message = required_err
	} else {
		err.Message = fmt.Sprintf("%s is required", fieldName)
	}

	return err
}

const (
	CodeMin          = "min"
	CodeMax          = "max"
	CodeLength       = "length"
	CodeEmail        = "email"
	CodeUUID         = "uuid"
	CodeURL          = "url"
	CodeEndsWith     = "ends-with"
	CodeStartsWith   = "starts-with"
	CodeAlpha        = "alpha"
	CodeNumeric      = "numeric"
	CodeAlphaNumeric = "alpha-numeric"
	CodeIsOneOf      = "is-one-of"
	CodeRefinement   = "refinement"
	CodeRequired     = "required"
	CodeContains     = "contains"
	CodeIs           = "is"
	CodeInvalidType  = "invalid-type"
)
