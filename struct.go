package validator

import "reflect"

type structAction[T any] struct {
	field          field
	refinement     func(T) error
	transformer    func(T) T
	refinementData RefinementData
}

type StructField[T any] struct {
	value         *T
	name          string
	optional      bool
	requiredError string
	actions       []structAction[T]
	abortEarly    bool
}

func (f *StructField[T]) addRefinement(fn func(T) error, refinementData RefinementData) {
	action := structAction[T]{refinement: fn, refinementData: refinementData}
	f.actions = append(f.actions, action)
}

func (f *StructField[T]) addTransformer(fn func(T) T) {
	action := structAction[T]{transformer: fn}
	f.actions = append(f.actions, action)
}

func (f *StructField[T]) _parse(errs *[]Error) bool {
	if f.value == nil {
		if !f.optional {
			*errs = append(*errs, requiredFieldErr(f.name, f.requiredError))
			return false
		}

		return true
	}

	if reflect.ValueOf(f.value).Elem().Kind() != reflect.Struct {
		*errs = append(*errs, Error{Field: f.name, Message: "value must be a struct", Code: CodeInvalidType})
		return false
	}

	isFieldParsedSuccessfully := true
	for _, action := range f.actions {
		ok := true
		if action.field != nil {
			ok = action.field._parse(errs)
		} else if action.refinement != nil {
			err := action.refinement(*f.value)
			if err != nil {
				ok = false
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

		if !ok {
			isFieldParsedSuccessfully = false
			if f.abortEarly {
				return false
			}
		}
	}

	return isFieldParsedSuccessfully
}

// AbortEarly stops the parsing of the field on the first error
func (f *StructField[T]) AbortEarly() *StructField[T] {
	f.abortEarly = true
	return f
}

// Optional makes the field optional
func (f *StructField[T]) Optional() *StructField[T] {
	f.optional = true
	return f
}

// Sets a custom error message if the field is missing
func (f *StructField[T]) RequiredError(message string) *StructField[T] {
	f.requiredError = message
	return f
}

// Fields take in fields of the struct and validates them
func (f *StructField[T]) Fields(fields ...field) *StructField[T] {
	for _, field := range fields {
		action := structAction[T]{field: field}
		f.actions = append(f.actions, action)
	}
	return f
}

// Refine lets you provide custom validation logic
func (f *StructField[T]) Refine(fn func(T) error, refinementData ...RefinementData) *StructField[T] {
	var newRefinementData RefinementData
	if len(refinementData) > 0 {
		newRefinementData = refinementData[0]
	}

	f.addRefinement(fn, newRefinementData)
	return f
}

// Transform "transforms" the field value.
func (f *StructField[T]) Transform(fn func(T) T) *StructField[T] {
	f.addTransformer(fn)
	return f
}

// Parse parses the field and returns a slice of Error.
func (f *StructField[T]) Parse() []Error {
	var errs []Error
	f._parse(&errs)
	return errs
}

// Struct takes a pointer to a struct and a variadic argument 'name'.
// Even if multiple values are passed for 'name', only the first value will be considered.
func Struct[T any](value *T, name ...string) *StructField[T] {
	field := StructField[T]{
		value: value,
	}

	if len(name) > 0 {
		field.name = name[0]
	}

	return &field
}
