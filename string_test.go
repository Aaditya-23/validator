package validator

import (
	"errors"
	"strings"
	"testing"
)

func TestStringMin(t *testing.T) {
	value := 5
	goodInput := "aaditya"
	badInput := "aadi"
	var errs []Error

	errs = String(&goodInput).Min(value).Parse()
	if len(errs) > 0 {
		t.Errorf("expected no error")
	}

	errs = String(&badInput).Min(value).Parse()
	if len(errs) == 0 {
		t.Fatal("expected error")
	}
}
func TestStringMax(t *testing.T) {
	value := 5
	goodInput := "aadi"
	badInput := "aaditya"
	var errs []Error

	errs = String(&goodInput).Max(value).Parse()
	if len(errs) > 0 {
		t.Errorf("expected no error")
	}

	errs = String(&badInput).Max(value).Parse()
	if len(errs) == 0 {
		t.Fatal("expected error")
	}
}
func TestStringLength(t *testing.T) {
	value := 5
	goodInput := "aadit"
	badInput := "aaditya"
	var errs []Error

	errs = String(&goodInput).Length(value).Parse()
	if len(errs) > 0 {
		t.Errorf("expected no error")
	}

	errs = String(&badInput).Length(value).Parse()
	if len(errs) == 0 {
		t.Fatal("expected error")
	}
}
func TestStringContains(t *testing.T) {
	value := "dit"
	goodInput := "aaditya"
	badInput := "aadi"
	var errs []Error

	errs = String(&goodInput).Contains(value).Parse()
	if len(errs) > 0 {
		t.Errorf("expected no error")
	}

	errs = String(&badInput).Contains(value).Parse()
	if len(errs) == 0 {
		t.Fatal("expected error")
	}
}

// func TestEmail(t *testing.T) {
// 	value := 5
// 	goodInput := "aadi"
// 	badInput := "aaditya"
// 	var errs []Error

// 	errs = String(&goodInput, field).Min(value).Parse()
// 	if len(errs) > 0 {
// 		t.Errorf("expected no error")
// 	}

// 	errs = String(&badInput, field).Min(value).Parse()
// 	if len(errs) == 0 {
// 		t.Fatal("expected error")
// 	}
// }
// func TestUUID(t *testing.T) {
// 	value := 5
// 	goodInput := "aadi"
// 	badInput := "aaditya"
// 	var errs []Error

// 	errs = String(&goodInput, field).Min(value).Parse()
// 	if len(errs) > 0 {
// 		t.Errorf("expected no error")
// 	}

// 	errs = String(&badInput, field).Min(value).Parse()
// 	if len(errs) == 0 {
// 		t.Fatal("expected error")
// 	}
// }
// func TestURL(t *testing.T) {
// 	value := 5
// 	goodInput := "aadi"
// 	badInput := "aaditya"
// 	var errs []Error

// 	errs = String(&goodInput, field).Min(value).Parse()
// 	if len(errs) > 0 {
// 		t.Errorf("expected no error")
// 	}

//		errs = String(&badInput, field).Min(value).Parse()
//		if len(errs) == 0 {
//			t.Fatal("expected error")
//		}
//	}
func TestStringEndsWith(t *testing.T) {
	value := "tya"
	goodInput := "aaditya"
	badInput := "aadi"
	var errs []Error

	errs = String(&goodInput).EndsWith(value).Parse()
	if len(errs) > 0 {
		t.Errorf("expected no error")
	}

	errs = String(&badInput).EndsWith(value).Parse()
	if len(errs) == 0 {
		t.Fatal("expected error")
	}
}
func TestStringStartsWith(t *testing.T) {
	value := "aadi"
	goodInput := "aaditya"
	badInput := "verma"
	var errs []Error

	errs = String(&goodInput).StartsWith(value).Parse()
	if len(errs) > 0 {
		t.Errorf("expected no error")
	}

	errs = String(&badInput).StartsWith(value).Parse()
	if len(errs) == 0 {
		t.Fatal("expected error")
	}
}
func TestStringAlpha(t *testing.T) {
	goodInput := "aaditya"
	badInput := "aadi@23"
	var errs []Error

	errs = String(&goodInput).Alpha().Parse()
	if len(errs) > 0 {
		t.Errorf("expected no error")
	}

	errs = String(&badInput).Alpha().Parse()
	if len(errs) == 0 {
		t.Fatal("expected error")
	}
}
func TestStringNumeric(t *testing.T) {
	goodInput := "9999999999"
	badInput := "99aa99"
	var errs []Error

	errs = String(&goodInput).Numeric().Parse()
	if len(errs) > 0 {
		t.Errorf("expected no error")
	}

	errs = String(&badInput).Numeric().Parse()
	if len(errs) == 0 {
		t.Fatal("expected error")
	}
}
func TestStringAlphaNumeric(t *testing.T) {
	goodInput := "aaditya23"
	badInput := "aadi@23"
	var errs []Error

	errs = String(&goodInput).AlphaNumeric().Parse()
	if len(errs) > 0 {
		t.Errorf("expected no error")
	}

	errs = String(&badInput).AlphaNumeric().Parse()
	if len(errs) == 0 {
		t.Fatal("expected error")
	}
}
func TestStringIsOneOf(t *testing.T) {
	values := []string{"aadi", "aaditya", "eddie"}
	goodInput := "aaditya"
	badInput := "aadi@23"
	var errs []Error

	errs = String(&goodInput).IsOneOf(values).Parse()
	if len(errs) > 0 {
		t.Errorf("expected no error")
	}

	errs = String(&badInput).IsOneOf(values).Parse()
	if len(errs) == 0 {
		t.Fatal("expected error")
	}
}
func TestStringTrimSpace(t *testing.T) {
	input := "  aa di  "
	output := "aa di"

	errs := String(&input).TrimSpace().Parse()
	if len(errs) > 0 {
		t.Errorf("expected no error")
	}

	if input != output {
		t.Errorf("input not trimmed properly")
	}
}

func TestStringToLowerCase(t *testing.T) {
	input := "AaDiTyA"
	output := "aaditya"

	errs := String(&input).ToLowerCase().Parse()
	if len(errs) > 0 {
		t.Errorf("expected no error")
	}

	if input != output {
		t.Errorf("input not converted to lower case properly")
	}
}

func TestStringRefine(t *testing.T) {
	palindrome := "madam"

	errs := String(&palindrome).
		Refine(func(name string) error {
			var reversed string
			for _, c := range palindrome {
				reversed = string(c) + reversed
			}

			if reversed != palindrome {
				return errors.New("not a palindrome")
			}

			return nil
		}).
		Parse()

	if len(errs) > 0 {
		t.Errorf("expected no error")
	}
}

func TestStringTransform(t *testing.T) {
	input := "aaditya-verma"
	output := "aaditya verma"

	errs := String(&input).Transform(func(s string) string {
		return strings.ReplaceAll(s, "-", " ")
	}).Parse()

	if len(errs) > 0 {
		t.Errorf("expected no error")
	}

	if input != output {
		t.Errorf("input not transformed properly")
	}
}
