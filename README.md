Validator is an input validation package.

## Features

- Provides a set of pre-defined validation rules for common use cases.
- Allows custom validation rules to be easily implemented.
- No reflection is used for validating the input.
- Granular error handling.

## Installation

You can install Validator by running:
```
`go get github.com/aaditya-23/validator`
```

## Usage

All the methods used will be executed in the same order as they are chained. The only exceptions are `AbortEarly` which stops the validation of the field on the first failed method and `Optional` which skips the validation of the field if it is empty.
All the validation methods takes in an optional error message to display when the validation fails.

#### Validating String

```
email := "me@mail.com"

// second argument if the field name and is optional. Even if you pass more than 1 value to the second argument, the first value will be used.
errs := m.String(&email, "email").Email().Parse()
if len(errs) > 0 {
    println(errs[0].Field, errs[0].Message, errs[0].Code)
}
```

#### Validating Number

```
age := 17
errs := m.Number(&age, "age").Min(18).Parse()
if len(errs) > 0 {
    println(errs[0].Field, errs[0].Message, errs[0].Code)
}
```

#### Validating Struct

```
type User struct {
    name string
    age  int
}

user := User{
    name: "John Wick",
    age: 59,
}

errs := m.Struct(&user, "user").
        Fields(
            m.String(&user.name, "name"),
            m.Number(&user.age, "age").Min(18),
        ).
        AbortEarly().
        Parse()

if len(errs) > 0 {
    println(errs[0].Field, errs[0].Message, errs[0].Code)
}
```

#### Validating Bool

```
userExists := false
errs := m.Bool(&userExists).Is(true).Parse()
if len(errs) > 0 {
    println(errs[0].Field, errs[0].Message, errs[0].Code)
}
```

#### Validating Slice

```
hobbies := []string{"programing", "programing", "programing"}

errs := m.Slice(&hobbies).Min(1).Parse()
if len(errs) > 0 {
    println(errs[0].Field, errs[0].Message, errs[0].Code)
}
```

#### Validating Maps

```
address := map[string]string{
		"Street": "123",
		"City":   "",
		"State":  "Virginia",
		"Zip":    "12345",
	},

errs := m.Map(email).Parse()
if len(errs) > 0 {
    println(errs[0].Field, errs[0].Message, errs[0].Code)
}
```

### Refinement

Refinement is a way to apply custom validation logic to the field.

```
type User struct {
    name string
    email  string
}

user := User{
    name: "John Wick",
    email: "me@mail.com",
}

errs := m.Struct(&user, "user").
        Fields(
            m.String(&user.name, "name"),
            m.String(&user.age, "age").Email(),
        ).
        Refine(func (u User) error {
            // custom validation logic
            return nil
        }).
        AbortEarly().
        Parse()

if len(errs) > 0 {
    println(errs[0].Field, errs[0].Message, errs[0].Code)
}
```

### Transformations

Transformations are a way to transform the field value

```
email := "me@mail.com"

errs := m.String(&email, "email").Email().ToLowerCase().Parse()
if len(errs) > 0 {
    println(errs[0].Field, errs[0].Message, errs[0].Code)
}
```
