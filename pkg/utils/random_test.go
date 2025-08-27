package utils

import (
	"testing"
)

func Test_RandomInt(t *testing.T) {
	min, max := 1, 10
	values := make(map[int]bool)
	for i := 0; i < 100; i++ {
		n := RandomInt(min, max)
		if n < min || n > max {
			t.Errorf("RandomInt(%d, %d) = %d; want between %d and %d", min, max, n, min, max)
		}
		values[n] = true
	}
	if len(values) < 2 {
		t.Errorf("RandomInt(%d, %d) did not produce multiple unique values", min, max)
	}

	// Edge case: min == max
	single := RandomInt(5, 5)
	if single != 5 {
		t.Errorf("RandomInt(5, 5) = %d; want 5", single)
	}
}
