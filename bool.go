package validator

import (
	"errors"
	"fmt"
)

type boolAction struct {
	validator      func() error
	refinement     func(bool) error
	transformer    func(bool) bool
	code           string
	refinementData RefinementData
}

type boolField struct {
	value         *bool
	name          string
	optional      bool
	requiredError string
	actions       []boolAction
	abortEarly    bool
}

func (f *boolField) addValidation(fn func() error, code string) {
	action := boolAction{validator: fn, code: code}
	f.actions = append(f.actions, action)
}

func (f *boolField) addRefinement(fn func(bool) error, refinementData RefinementData) {
	action := boolAction{refinement: fn, refinementData: refinementData}
	f.actions = append(f.actions, action)
}

func (f *boolField) addTransformer(fn func(bool) bool) {
	action := boolAction{transformer: fn}
	f.actions = append(f.actions, action)
}

func (f *boolField) _parse(errs *[]Error) bool {
	if f.value == nil {
		if !f.optional {
			*errs = append(*errs, requiredFieldErr(f.name, f.requiredError))
			return false
		}

		return true
	}

	isFieldParsedSuccessfully := true
	for _, action := range f.actions {
		isActionParsedSuccessfully := true

		if action.validator != nil {
			err := action.validator()
			if err != nil {
				isActionParsedSuccessfully = false
				*errs = append(*errs, Error{Field: f.name, Message: err.Error(), Code: action.code})
			}
		} else if action.refinement != nil {
			err := action.refinement(*f.value)
			if err != nil {
				isActionParsedSuccessfully = false
				me := Error{Field: f.name, Message: err.Error(), Code: CodeRefinement}
				if action.refinementData.Field != "" {
					me.Field = action.refinementData.Field
				}
				if action.refinementData.Code != "" {
					me.Code = action.refinementData.Code
				}

				*errs = append(*errs, me)
			}
		} else if action.transformer != nil {
			*f.value = action.transformer(*f.value)
			continue
		}

		if !isActionParsedSuccessfully {
			isFieldParsedSuccessfully = false
			if f.abortEarly {
				return false
			}
		}
	}

	return isFieldParsedSuccessfully
}

// AbortEarly stops the parsing of the field on the first error
func (f *boolField) AbortEarly() *boolField {
	f.abortEarly = true
	return f
}

// Optional makes the field optional
func (f *boolField) Optional() *boolField {
	f.optional = true
	return f
}

// Sets a custom error message if the field is missing
func (f *boolField) RequiredError(message string) *boolField {
	f.requiredError = message
	return f
}

// Is checks if the field value is equal to the provided boolean value
func (f *boolField) Is(value bool, message ...string) *boolField {
	fv := *f.value
	code := CodeIs

	validator := func() error {
		if value != fv {
			var msg string
			if len(message) > 0 {
				msg = message[0]
			} else {
				msg = fmt.Sprintf("%s should be %t", f.name, value)
			}

			return errors.New(msg)
		}
		return nil
	}

	f.addValidation(validator, code)
	return f
}

// Refine lets you provide custom validation logic
func (f *boolField) Refine(fn func(bool) error, refinementData ...RefinementData) *boolField {
	var newRefinementData RefinementData
	if len(refinementData) > 0 {
		newRefinementData = refinementData[0]
	}

	f.addRefinement(fn, newRefinementData)
	return f
}

// Transform "transforms" the field value.
func (f *boolField) Transform(fn func(bool) bool) *boolField {
	f.addTransformer(fn)
	return f
}

// Parse parses the field and returns a slice of Error.
func (f *boolField) Parse() []Error {
	var errs []Error
	f._parse(&errs)
	return errs
}

// Bool takes a pointer to a bool and a variadic argument 'name'.
// Even if multiple values are passed for 'name', only the first value will be considered.
func Bool(value *bool, name ...string) *boolField {
	field := boolField{
		value: value,
	}

	if len(name) > 0 {
		field.name = name[0]
	}

	return &field
}
