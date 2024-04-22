package validator

import (
	"errors"
	"fmt"
)

type mapAction[T comparable, K any] struct {
	validator      func() error
	refinement     func(map[T]K) error
	transformer    func(map[T]K)
	code           string
	refinementData RefinementData
}

type MapField[T comparable, K any] struct {
	value         *map[T]K
	name          string
	optional      bool
	requiredError string
	actions       []mapAction[T, K]
	abortEarly    bool
}

func (f *MapField[T, K]) addValidation(fn func() error, code string) {
	r := mapAction[T, K]{validator: fn, code: code}
	f.actions = append(f.actions, r)
}

func (f *MapField[T, K]) addRefinement(fn func(map[T]K) error, refinementData RefinementData) {
	r := mapAction[T, K]{refinement: fn, refinementData: refinementData}
	f.actions = append(f.actions, r)
}

func (f *MapField[T, K]) addTransformer(fn func(map[T]K)) {
	r := mapAction[T, K]{transformer: fn}
	f.actions = append(f.actions, r)
}

func (f *MapField[T, K]) _parse(errs *[]Error) bool {
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
			action.transformer(*f.value)
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
func (f *MapField[T, K]) AbortEarly() *MapField[T, K] {
	f.abortEarly = true
	return f
}

// Optional makes the field optional
func (f *MapField[T, K]) Optional() *MapField[T, K] {
	f.optional = true
	return f
}

// Sets a custom error message if the field is missing
func (f *MapField[T, K]) RequiredError(message string) *MapField[T, K] {
	f.requiredError = message
	return f
}

// Min sets the minimum number of entries the map should have.
func (f *MapField[T, K]) Min(size int, message ...string) *MapField[T, K] {
	fv := *f.value
	code := CodeMin

	validator := func() error {
		if len(fv) < size {
			var msg string
			if len(message) > 0 {
				msg = message[0]
			} else {
				msg = fmt.Sprintf("%s should have atleast %d entries", f.name, size)
			}

			return errors.New(msg)
		}

		return nil
	}

	f.addValidation(validator, code)
	return f
}

// Max sets the maximum number of entries for the map
func (f *MapField[T, K]) Max(size int, message ...string) *MapField[T, K] {
	fv := *f.value
	code := CodeMax

	validator := func() error {
		if len(fv) > size {
			var msg string
			if len(message) > 0 {
				msg = message[0]
			} else {
				msg = fmt.Sprintf("%s can have atmost %d entries", f.name, size)
			}

			return errors.New(msg)
		}

		return nil
	}

	f.addValidation(validator, code)
	return f
}

// Refine lets you provide custom validation logic
func (f *MapField[T, K]) Refine(fn func(map[T]K) error, refinementData ...RefinementData) *MapField[T, K] {
	var newRefinementData RefinementData
	if len(refinementData) > 0 {
		newRefinementData = refinementData[0]
	}

	f.addRefinement(fn, newRefinementData)
	return f
}

// Transform "transforms" the field value.
func (f *MapField[T, K]) Transform(fn func(map[T]K)) *MapField[T, K] {
	f.addTransformer(fn)
	return f
}

// Parse parses the field and returns a slice of Error.
func (f *MapField[T, K]) Parse() []Error {
	var errs []Error
	f._parse(&errs)
	return errs
}

// Map takes a pointer to a map and a variadic argument 'name'.
// Even if multiple values are passed for 'name', only the first value will be considered.
func Map[T comparable, K any](value *map[T]K, name ...string) *MapField[T, K] {
	field := MapField[T, K]{
		value: value,
	}

	if len(name) > 0 {
		field.name = name[0]
	}

	return &field
}
