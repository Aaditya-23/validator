package validator

import (
	"errors"
	"testing"
)

func TestMapMin(t *testing.T) {
	value := 3
	goodInput := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	badInput := map[string]int{"a": 1, "b": 2}
	var errs []Error

	errs = Map(&goodInput).
		Min(value).
		Parse()

	if len(errs) > 0 {
		t.Error("expected no error")
	}

	errs = Map(&badInput).
		Min(value).
		Parse()

	if len(errs) == 0 {
		t.Error("expected error")
	}
}

func TestMapMax(t *testing.T) {
	value := 3
	goodInput := map[string]int{"a": 1, "b": 2}
	badInput := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	var errs []Error

	errs = Map(&goodInput).
		Max(value).
		Parse()

	if len(errs) > 0 {
		t.Error("expected no error")
	}

	errs = Map(&badInput).
		Max(value).
		Parse()

	if len(errs) == 0 {
		t.Error("expected error")
	}
}

func TestMapRefine(t *testing.T) {
	input := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}

	errs := Map(&input).
		Refine(func(m map[string]int) error {
			_, ok := m["c"]
			if !ok {
				return errors.New("key 'c' not found")
			}

			return nil
		}).
		Parse()

	if len(errs) > 0 {
		t.Error("expected no error")
	}
}

func TestMapTransform(t *testing.T) {
	input := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	output := map[string]int{"a": 2, "b": 4, "c": 6, "d": 8}

	errs := Map(&input).
		Transform(func(m map[string]int) {
			for k := range m {
				m[k] = m[k] * 2
			}
		}).
		Parse()

	if len(errs) > 0 {
		t.Error("expected no error")
	}

	for k := range input {
		if input[k] != output[k] {
			t.Error("input not transformed properly")
			break
		}
	}
}
