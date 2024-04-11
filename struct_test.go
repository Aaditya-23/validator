package validator

import (
	"errors"
	"strings"
	"testing"
)

func TestStruct(t *testing.T) {

	type User struct {
		name string
		age  int
	}

	user := User{
		name: "aaditya",
		age:  21,
	}

	errs := Struct(&user).Fields(
		String(&user.name).Alpha(),
		Number(&user.age).Min(18),
	).Parse()

	if len(errs) > 0 {
		t.Error("expected no error")
	}
}

func TestStructRefine(t *testing.T) {

	type User struct {
		name          string
		age           int
		palindromeKey string
	}

	user := User{
		name:          "aaditya",
		age:           21,
		palindromeKey: "123321",
	}

	errs := Struct(&user).
		Fields(
			String(&user.name).Alpha(),
			Number(&user.age).Min(18),
		).
		Refine(func(u User) error {
			var reversed string
			for _, c := range user.palindromeKey {
				reversed = string(c) + reversed
			}

			if reversed != user.palindromeKey {
				return errors.New("not a palindrome")
			}

			return nil
		}).
		Parse()

	if len(errs) > 0 {
		t.Error("expected no error")
	}
}

func TestStructTransform(t *testing.T) {

	type User struct {
		name, email, domainName string
	}

	user := User{
		name:  "aaditya",
		email: "aaditya220055@gmail.com",
	}

	errs := Struct(&user).
		Fields(
			String(&user.name).Alpha(),
			String(&user.email).Email(),
		).
		Transform(func(u User) User {
			domain := strings.Split(u.email, "@")[1]
			u.domainName = strings.Split(domain, ".")[0]

			return u
		}).
		Parse()

	if len(errs) > 0 {
		t.Error("expected no error")
	}

	if user.domainName != "gmail" {
		t.Error("input not transformed properly")
	}
}
