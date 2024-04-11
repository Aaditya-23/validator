package validator

import (
	"errors"
	"fmt"
)

type number interface {
	int8 | int16 | int | int32 | int64 | float32 | float64
}

type numberAction[T number] struct {
	validator      func() error
	refinement     func(T) error
	transformer    func(T) T
	code           string
	refinementData RefinementData
}

type numberField[T number] struct {
	value         *T
	name          string
	optional      bool
	requiredError string
	actions       []numberAction[T]
	abortEarly    bool
}

func (f *numberField[T]) addValidation(fn func() error, code string) {
	action := numberAction[T]{validator: fn, code: code}
	f.actions = append(f.actions, action)
}

func (f *numberField[T]) addRefinement(fn func(T) error, refinementData RefinementData) {
	action := numberAction[T]{refinement: fn, refinementData: refinementData}
	f.actions = append(f.actions, action)
}

func (f *numberField[T]) addTransformer(fn func(T) T) {
	action := numberAction[T]{transformer: fn}
	f.actions = append(f.actions, action)
}

func (f *numberField[T]) _parse(errs *[]Error) bool {
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
func (f *numberField[T]) AbortEarly() *numberField[T] {
	f.abortEarly = true
	return f
}

// Optional makes the field optional
func (f *numberField[T]) Optional() *numberField[T] {
	f.optional = true
	return f
}

// Sets a custom error message if the field is missing
func (f *numberField[T]) RequiredError(message string) *numberField[T] {
	f.requiredError = message
	return f
}

// Min sets the minimum value for the field.
func (f *numberField[T]) Min(value T, message ...string) *numberField[T] {
	fv := *f.value
	code := CodeMin

	validator := func() error {
		if fv < value {
			var msg string
			if len(message) > 0 {
				msg = message[0]
			} else {
				msg = fmt.Sprintf("%s must be atleast %v", f.name, value)
			}

			return errors.New(msg)
		}

		return nil
	}

	f.addValidation(validator, code)
	return f
}

// Max sets the maximum value for the field.
func (f *numberField[T]) Max(value T, message ...string) *numberField[T] {
	fv := *f.value
	code := CodeMax

	validator := func() error {
		if fv > value {
			var msg string
			if len(message) > 0 {
				msg = message[0]
			} else {
				msg = fmt.Sprintf("%s can be atmost %v", f.name, value)
			}

			return errors.New(msg)
		}

		return nil
	}

	f.addValidation(validator, code)
	return f
}

// Refine lets you provide custom validation logic
func (f *numberField[T]) Refine(fn func(T) error, refinementData ...RefinementData) *numberField[T] {
	var newRefinementData RefinementData
	if len(refinementData) > 0 {
		newRefinementData = refinementData[0]
	}

	f.addRefinement(fn, newRefinementData)
	return f
}

// Transform "transforms" the field value.
func (f *numberField[T]) Transform(fn func(T) T) *numberField[T] {
	f.addTransformer(fn)
	return f
}

// Parse parses the field and returns a slice of Error.
func (f *numberField[T]) Parse() []Error {
	var errs []Error
	f._parse(&errs)
	return errs
}

// Number takes a pointer to a number and a variadic argument 'name'.
// Even if multiple values are passed for 'name', only the first value will be considered.
func Number[T number](value *T, name ...string) *numberField[T] {
	field := numberField[T]{
		value: value,
	}

	if len(name) > 0 {
		field.name = name[0]
	}

	return &field
}
