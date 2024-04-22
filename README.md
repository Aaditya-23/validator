Validator is an input validation package.

## Features

- Provides a set of pre-defined validation rules for common use cases.
- Easily add custom validation logic.
- Granular error handling.

## Installation

You can install Validator by running:

```
`go get github.com/aaditya-23/validator`
```

## Usage

All the methods used will be executed in the same order as they are chained. The only exceptions are `AbortEarly` which stops the validation of the field on the first failed method and `Optional` which skips the validation of the field if it is empty.
All the validation methods takes in an optional error message to display when the validation fails.

### Importing the package

```go
import v "github.com/aaditya-23/validator"
```

#### Validating String

```go
email := "me@mail.com"

// second argument is the field name and is optional. Even if you pass more than 1 value to the second argument, the first value will be used.
errs := v.String(&email, "email").
        Email().
        Parse()
if len(errs) > 0 {
    println(errs[0].Field, errs[0].Message, errs[0].Code)
}
```

### Custom Error messages

You can pass an optional error message to the rule

```go
email := "me@mail.com"
errs := v.String(&email, "email").
        Email("please enter a valid email").
        Parse()
```

#### Validating Number

```go
age := 17
errs := v.Number(&age, "age").
        Min(18, "only adults can access this content").
        Parse()
if len(errs) > 0 {
    println(errs[0].Field, errs[0].Message, errs[0].Code)
}
```

#### Validating Struct

```go
type User struct {
    name string
    age  int
}

user := User{
    name: "John Wick",
    age: 59,
}

errs := v.Struct(&user, "user").
        Fields(
            v.String(&user.name, "name"),
            v.Number(&user.age, "age").Min(18),
        ).
        AbortEarly().
        Parse()

if len(errs) > 0 {
    println(errs[0].Field, errs[0].Message, errs[0].Code)
}
```

#### Validating Bool

```go
userExists := false
errs := v.Bool(&userExists).
        Is(true).
        Parse()

if len(errs) > 0 {
    println(errs[0].Field, errs[0].Message, errs[0].Code)
}
```

#### Validating Slice

```go
hobbies := []string{"programing", "programing", "programing"}

errs := v.Slice(&hobbies, "hobbies").
        Min(1).
        Parse()

if len(errs) > 0 {
    println(errs[0].Field, errs[0].Message, errs[0].Code)
}
```

#### Validating Maps

```go
address := map[string]string{
		"Street": "123",
		"City":   "",
		"State":  "Virginia",
		"Zip":    "12345",
	},

errs := v.Map(&address)
        .Parse()

if len(errs) > 0 {
    println(errs[0].Field, errs[0].Message, errs[0].Code)
}
```

### Refinement

Refinement is a way to apply custom validation logic to the field.

```go
type User struct {
    name string
    email  string
}

user := User{
    name: "John Wick",
    email: "me@mail.com",
}

errs := v.Struct(&user, "user").
        Fields(
            v.String(&user.name, "name"),
            v.String(&user.email, "email").Email(),
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

When a refinement fails, the default Error 'Code' will be 'refinement' and the 'Field' will be the name of the field it is chained to. However, when working with structs, you may want to provide a different 'Field' name or 'Code'. You can do this by passing `RefinementData` as the second argument.

```go
type User struct {
    name     string
    email    string
    username string
}

user := User {
    name:     "John Wick",
    email:    "me@mail.com",
    username: "john@wick"
}

errs := v.Struct(&user, "user").
        Fields(
            v.String(&user.name, "name"),
            v.String(&user.email, "email").Email(),
            v.String(&user.username, "username"),
        ).
        Refine(func (u User) error {
            if (strings.Contains(u.username, "@")){
                return errors.New("username should not contain '@'")
            }
            
            return nil
        }, v.RefinementData {Field: "username", Code: "invalid-value"}).
        AbortEarly().
        Parse()
```

### Transformations

Transformations are a way to transform the field value

```go
email := "me@mail.com"

errs := m.String(&email, "email").
        ToLowerCase().
        Email().
        Parse()

if len(errs) > 0 {
    println(errs[0].Field, errs[0].Message, errs[0].Code)
}
```
