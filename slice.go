package validator

import (
	"errors"
	"fmt"
)

type sliceAction[T any] struct {
	validator      func() error
	refinement     func([]T) error
	transformer    func([]T) []T
	code           string
	refinementData RefinementData
}

type SliceField[T any] struct {
	value         *[]T
	name          string
	optional      bool
	requiredError string
	actions       []sliceAction[T]
	abortEarly    bool
}

func (f *SliceField[T]) addValidation(fn func() error, code string) {
	r := sliceAction[T]{validator: fn, code: code}
	f.actions = append(f.actions, r)
}

func (f *SliceField[T]) addRefinement(fn func([]T) error, refinementData RefinementData) {
	r := sliceAction[T]{refinement: fn, refinementData: refinementData}
	f.actions = append(f.actions, r)
}

func (f *SliceField[T]) addTransformer(fn func([]T) []T) {
	r := sliceAction[T]{transformer: fn}
	f.actions = append(f.actions, r)
}

func (f *SliceField[T]) _parse(errs *[]Error) bool {
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
func (f *SliceField[T]) AbortEarly() *SliceField[T] {
	f.abortEarly = true
	return f
}

// Optional makes the field optional
func (f *SliceField[T]) Optional() *SliceField[T] {
	f.optional = true
	return f
}

// Sets a custom error message if the field is missing
func (f *SliceField[T]) RequiredError(message string) *SliceField[T] {
	f.requiredError = message
	return f
}

// Min sets the minimum length of the slice
func (f *SliceField[T]) Min(length int, message ...string) *SliceField[T] {
	fv := *f.value
	code := CodeMin

	validator := func() error {
		if len(fv) < length {
			var msg string
			if len(message) > 0 {
				msg = message[0]
			} else {
				msg = fmt.Sprintf("%s must have atleast %d items", f.name, length)
			}

			return errors.New(msg)
		}
		return nil
	}

	f.addValidation(validator, code)
	return f
}

// Max sets the maximum length of the slice
func (f *SliceField[T]) Max(length int, message ...string) *SliceField[T] {
	fv := *f.value
	rule := "max"

	validator := func() error {
		if len(fv) > length {
			var msg string
			if len(message) > 0 {
				msg = message[0]
			} else {
				msg = fmt.Sprintf("%s must have atmost %d items", f.name, length)
			}

			return errors.New(msg)
		}
		return nil
	}

	f.addValidation(validator, rule)
	return f
}

// Length checks if the slice has exactly the provided length
func (f *SliceField[T]) Length(value int, message ...string) *SliceField[T] {
	fv := *f.value
	rule := "length"

	validator := func() error {
		if len(fv) != value {
			var msg string
			if len(message) > 0 {
				msg = message[0]
			} else {
				msg = fmt.Sprintf("%s must have %d items", f.name, value)
			}

			return errors.New(msg)
		}
		return nil
	}

	f.addValidation(validator, rule)
	return f
}

// Refine lets you provide custom validation logic
func (f *SliceField[T]) Refine(fn func([]T) error, refinementData ...RefinementData) *SliceField[T] {
	var newRefinementData RefinementData
	if len(refinementData) > 0 {
		newRefinementData = refinementData[0]
	}

	f.addRefinement(fn, newRefinementData)
	return f
}

// Transform "transforms" the field value.
func (f *SliceField[T]) Transform(fn func([]T) []T) *SliceField[T] {
	f.addTransformer(fn)
	return f
}

// Parse parses the field and returns a slice of Error.
func (f *SliceField[T]) Parse() []Error {
	var errs []Error
	f._parse(&errs)
	return errs
}

// Slice takes a pointer to a slice and a variadic argument 'name'.
// Even if multiple values are passed for 'name', only the first value will be considered.
func Slice[T any](value *[]T, name ...string) *SliceField[T] {
	field := SliceField[T]{
		value: value,
	}

	if len(name) > 0 {
		field.name = name[0]
	}

	return &field
}
