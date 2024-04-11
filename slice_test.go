package validator

import (
	"errors"
	"testing"
)

func TestSliceMin(t *testing.T) {
	value := 3
	goodInput := []int{1, 2, 3, 4}
	badInput := []int{1, 2}
	var errs []Error

	errs = Slice(&goodInput).Min(value).Parse()
	if len(errs) > 0 {
		t.Error("expected no error")
	}

	errs = Slice(&badInput).Min(value).Parse()
	if len(errs) == 0 {
		t.Error("expected error")
	}
}

func TestSliceMax(t *testing.T) {
	value := 3
	goodInput := []int{1, 2, 3}
	badInput := []int{1, 2, 3, 4}
	var errs []Error

	errs = Slice(&goodInput).Max(value).Parse()
	if len(errs) > 0 {
		t.Error("expected no error")
	}

	errs = Slice(&badInput).Max(value).Parse()
	if len(errs) == 0 {
		t.Error("expected error")
	}
}

func TestSliceLength(t *testing.T) {
	value := 3
	goodInput := []int{1, 2, 3}
	badInput := []int{1, 2}
	var errs []Error

	errs = Slice(&goodInput).Length(value).Parse()
	if len(errs) > 0 {
		t.Error("expected no error")
	}

	errs = Slice(&badInput).Length(value).Parse()
	if len(errs) == 0 {
		t.Error("expected error")
	}
}

func TestSliceRefine(t *testing.T) {

	input := []int{1, 2, 3, 4}

	errs := Slice(&input).
		Refine(func(v []int) error {
			for i := range input {
				if i > 0 && input[i] < input[i-1] {
					return errors.New("must be sorted")
				}
			}
			return nil
		}).
		Parse()

	if len(errs) > 0 {
		t.Error("expected no error")
	}
}

func TestSliceTransform(t *testing.T) {

	input := []int{1, 2, 3, 4}
	output := []int{2, 4, 6, 8}

	errs := Slice(&input).
		Transform(func(v []int) []int {
			for i := range v {
				v[i] = v[i] * 2
			}

			return v
		}).
		Parse()

	if len(errs) > 0 {
		t.Error("expected no error")
	}

	for i := range input {
		if input[i] != output[i] {
			t.Errorf("input not transformed properly")
			break
		}
	}

}
