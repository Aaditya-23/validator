package validator

import (
	"errors"
	"testing"
)

func TestNumberMin(t *testing.T) {
	value := 5
	goodInput := 7
	badInput := 3
	var errs []Error

	errs = Number(&goodInput).Min(value).Parse()
	if len(errs) > 0 {
		t.Error("expected no error")
	}

	errs = Number(&badInput).Min(value).Parse()
	if len(errs) == 0 {
		t.Error("expected error")
	}
}

func TestNumberMax(t *testing.T) {
	value := 5
	goodInput := 3
	badInput := 7
	var errs []Error

	errs = Number(&goodInput).Max(value).Parse()
	if len(errs) > 0 {
		t.Error("expected no error")
	}

	errs = Number(&badInput).Max(value).Parse()
	if len(errs) == 0 {
		t.Error("expected error")
	}
}

func TestNumberRefine(t *testing.T) {
	input := 5

	errs := Number(&input).
		Refine(func(i int) error {
			if i%2 == 0 {
				return nil
			}
			return errors.New("not an even number")
		}).
		Parse()

	if len(errs) == 0 {
		t.Error("expected error")
	}
}

func TestNumberTransform(t *testing.T) {
	input := 5
	output := 10

	errs := Number(&input).
		Transform(func(i int) int {
			return i * 2
		}).
		Parse()

	if len(errs) > 0 {
		t.Error("expected no error")
	}

	if input != output {
		t.Error("input not transformed properly")
	}
}
