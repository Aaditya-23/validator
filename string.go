package validator

import (
	"errors"
	"fmt"
	"strings"
)

type stringAction struct {
	validator      func() error
	refinement     func(string) error
	transformer    func(string) string
	code           string
	refinementData RefinementData
}

type StringField struct {
	value         *string
	name          string
	optional      bool
	requiredError string
	actions       []stringAction
	abortEarly    bool
}

func (f *StringField) addValidation(fn func() error, code string) {
	action := stringAction{validator: fn, code: code}
	f.actions = append(f.actions, action)
}

func (f *StringField) addRefinement(fn func(string) error, refinementData RefinementData) {
	action := stringAction{refinement: fn, refinementData: refinementData}
	f.actions = append(f.actions, action)
}

func (f *StringField) addTransformer(fn func(string) string) {
	action := stringAction{transformer: fn}
	f.actions = append(f.actions, action)
}

func (f *StringField) _parse(errs *[]Error) bool {
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
func (f *StringField) AbortEarly() *StringField {
	f.abortEarly = true
	return f
}

// Optional makes the field optional
func (f *StringField) Optional() *StringField {
	f.optional = true
	return f
}

// Sets a custom error message if the field is missing
func (f *StringField) RequiredError(message string) *StringField {
	f.requiredError = message
	return f
}

// Min checks if the field value has the provided minimum length
func (f *StringField) Min(length int, message ...string) *StringField {
	fv := *f.value
	code := CodeMin

	validator := func() error {
		if len(fv) < length {
			var msg string
			if len(message) > 0 {
				msg = message[0]
			} else {
				msg = fmt.Sprintf("%s should have atleast %d characters", f.name, length)
			}

			return errors.New(msg)
		}

		return nil
	}

	f.addValidation(validator, code)
	return f
}

// Max checks if the field value has the provided maximum length
func (f *StringField) Max(length int, message ...string) *StringField {
	fv := *f.value
	code := CodeMax

	validator := func() error {
		if len(fv) > length {
			var msg string
			if len(message) > 0 {
				msg = message[0]
			} else {
				msg = fmt.Sprintf("%s can have atmost %d characters", f.name, length)
			}

			return errors.New(msg)
		}

		return nil
	}

	f.addValidation(validator, code)
	return f
}

// Length checks if the field value has the provided length
func (f *StringField) Length(value int, message ...string) *StringField {
	fv := *f.value
	code := CodeLength

	validator := func() error {
		if len(fv) != value {
			var msg string
			if len(message) > 0 {
				msg = message[0]
			} else {
				msg = fmt.Sprintf("%s should have %d characters", f.name, value)
			}

			return errors.New(msg)
		}

		return nil
	}

	f.addValidation(validator, code)
	return f
}

// Contains checks if the field value contains the provided substring
func (f *StringField) Contains(substr string, message ...string) *StringField {
	fv := *f.value
	code := CodeContains

	validator := func() error {
		if !strings.Contains(fv, substr) {
			var msg string
			if len(message) > 0 {
				msg = message[0]
			} else {
				msg = fmt.Sprintf("%s should contain %s", f.name, substr)
			}

			return errors.New(msg)
		}

		return nil
	}

	f.addValidation(validator, code)
	return f
}

// Email checks if the field value is a valid email address
func (f *StringField) Email(message ...string) *StringField {
	fv := *f.value
	code := CodeEmail

	validator := func() error {
		isEmail := emailRegex.MatchString(fv)
		if !isEmail {
			var msg string
			if len(message) > 0 {
				msg = message[0]
			} else {
				msg = fmt.Sprintf("%s is not a valid email", f.name)
			}

			return errors.New(msg)
		}

		return nil
	}

	f.addValidation(validator, code)
	return f
}

// UUID checks if the field value is a valid UUID
func (f *StringField) UUID(message ...string) *StringField {
	fv := *f.value
	code := CodeUUID

	validator := func() error {
		isEmail := uuidRegex.MatchString(fv)
		if !isEmail {
			var msg string
			if len(message) > 0 {
				msg = message[0]
			} else {
				msg = fmt.Sprintf("%s is not a valid UUID", f.name)
			}

			return errors.New(msg)
		}

		return nil
	}

	f.addValidation(validator, code)
	return f
}

// URL checks if the field value is a valid URL
func (f *StringField) URL(message ...string) *StringField {
	fv := *f.value
	code := CodeURL

	validator := func() error {
		isURL := urlRegex.MatchString(fv)
		if !isURL {
			var msg string
			if len(message) > 0 {
				msg = message[0]
			} else {
				msg = fmt.Sprintf("%s is not a valid URL", f.name)
			}

			return errors.New(msg)
		}

		return nil
	}

	f.addValidation(validator, code)
	return f
}

// EndsWith checks if the field value ends with the provided value
func (f *StringField) EndsWith(value string, message ...string) *StringField {
	fv := *f.value
	code := CodeEndsWith

	validator := func() error {
		if !strings.HasSuffix(fv, value) {
			var msg string
			if len(message) > 0 {
				msg = message[0]
			} else {
				msg = fmt.Sprintf("%s does not ends with %s", f.name, value)
			}

			return errors.New(msg)
		}
		return nil
	}

	f.addValidation(validator, code)
	return f
}

// StartsWith checks if the field value starts with the provided value
func (f *StringField) StartsWith(value string, message ...string) *StringField {
	fv := *f.value
	code := CodeStartsWith

	validator := func() error {
		if !strings.HasPrefix(fv, value) {
			var msg string
			if len(message) > 0 {
				msg = message[0]
			} else {
				msg = fmt.Sprintf("%s does not starts with %s", f.name, value)
			}

			return errors.New(msg)
		}
		return nil
	}

	f.addValidation(validator, code)
	return f
}

// Alpha checks if the field value contains only alphabets
func (f *StringField) Alpha(message ...string) *StringField {
	fv := *f.value
	code := CodeAlpha

	validator := func() error {
		isAlpha := alphaRegex.MatchString(fv)
		if !isAlpha {
			var msg string
			if len(message) > 0 {
				msg = message[0]
			} else {
				msg = fmt.Sprintf("%s should contain only alphabets", f.name)
			}

			return errors.New(msg)
		}

		return nil
	}

	f.addValidation(validator, code)
	return f
}

// Numeric checks if the field value contains only numbers
func (f *StringField) Numeric(message ...string) *StringField {
	fv := *f.value
	code := "numeric"

	validator := func() error {
		isNumeric := numericRegex.MatchString(fv)
		if !isNumeric {
			var msg string
			if len(message) > 0 {
				msg = message[0]
			} else {
				msg = fmt.Sprintf("%s should contain only numbers", f.name)
			}

			return errors.New(msg)
		}

		return nil
	}

	f.addValidation(validator, code)
	return f
}

// AlphaNumeric checks if the field value contains only alphabets and numbers
func (f *StringField) AlphaNumeric(message ...string) *StringField {
	fv := *f.value
	code := CodeAlphaNumeric

	validator := func() error {
		isAlphaNumeric := alphaNumericRegex.MatchString(fv)
		if !isAlphaNumeric {
			var msg string
			if len(message) > 0 {
				msg = message[0]
			} else {
				msg = fmt.Sprintf("%s should contain only alphabets and numbers", f.name)
			}

			return errors.New(msg)
		}

		return nil
	}

	f.addValidation(validator, code)
	return f
}

// IsOneOf checks if the field value is one of the values passed in the slice
func (f *StringField) IsOneOf(values []string, message ...string) *StringField {
	fv := *f.value
	code := CodeIsOneOf

	validator := func() error {
		for _, value := range values {
			if value == fv {
				return nil
			}
		}

		var msg string
		if len(message) > 0 {
			msg = message[0]
		} else {
			msg = fmt.Sprintf("%s can only be %s", f.name, strings.Join(values, ", "))
		}

		return errors.New(msg)
	}

	f.addValidation(validator, code)
	return f
}

// TrimSpace trims the leading and trailing spaces from the field value
func (f *StringField) TrimSpace() *StringField {
	fn := func(value string) string {
		return strings.TrimSpace(value)
	}

	f.addTransformer(fn)
	return f
}

// ToUpperCase converts the lowercase characters to uppercase
func (f *StringField) ToLowerCase() *StringField {
	fn := func(value string) string {
		return strings.ToLower(value)
	}

	f.addTransformer(fn)
	return f
}

// Refine lets you provide custom validation logic
func (f *StringField) Refine(fn func(field string) error, refinementData ...RefinementData) *StringField {
	var newRefinementData RefinementData
	if len(refinementData) > 0 {
		newRefinementData = refinementData[0]
	}

	f.addRefinement(fn, newRefinementData)
	return f
}

// Transform "transforms" the field value.
func (f *StringField) Transform(fn func(string) string) *StringField {
	f.addTransformer(fn)
	return f
}

// Parse parses the field and returns a slice of Error.
func (f *StringField) Parse() []Error {
	var errs []Error
	f._parse(&errs)
	return errs
}

// String takes a pointer to a string and a variadic argument 'name'.
// Even if multiple values are passed for 'name', only the first value will be considered.
func String(value *string, name ...string) *StringField {
	field := StringField{
		value: value,
	}

	if len(name) > 0 {
		field.name = name[0]
	}

	return &field
}
