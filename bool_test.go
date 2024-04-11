package validator

import (
	"errors"
	"testing"
)

func TestBoolIs(t *testing.T) {
	value := true
	goodInput := true
	badInput := false
	var errs []Error

	errs = Bool(&goodInput).Is(value).Parse()
	if len(errs) > 0 {
		t.Error("expected no error")
	}

	errs = Bool(&badInput).Is(value).Parse()
	if len(errs) == 0 {
		t.Error("expected error")
	}
}

func TestBoolRefine(t *testing.T) {
	input := true
	dbCall := func() bool {
		return true
	}

	errs := Bool(&input).
		Refine(func(b bool) error {
			// make a fake db call to check if user exists
			exists := dbCall()
			if !exists {
				return errors.New("user does not exist")
			}

			return nil
		}).Parse()

	if len(errs) > 0 {
		t.Error("expected no error")
	}
}

func TestBoolTransform(t *testing.T) {
	input := true
	invert := func(b bool) bool {
		return !b
	}

	errs := Bool(&input).
		Transform(invert).
		Parse()

	if len(errs) > 0 {
		t.Error("expected no error")
	}

	if input == true {
		t.Error("input not transformed properly")
	}
}
