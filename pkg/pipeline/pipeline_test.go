package pipeline

import (
	"fmt"
	"testing"
)

// sliceSource is a test Source backed by a slice.
type sliceSource[T any] struct {
	items []T
	index int
}

func newSliceSource[T any](items []T) *sliceSource[T] {
	return &sliceSource[T]{items: items}
}

func (s *sliceSource[T]) Next() (T, bool, error) {
	var zero T
	if s.index >= len(s.items) {
		return zero, false, nil
	}
	val := s.items[s.index]
	s.index++
	return val, true, nil
}

func TestMapFilter(t *testing.T) {
	src := newSliceSource([]int{1, 2, 3, 4, 5})
	doubled := NewMapFilter(src, func(n int) (int, bool, error) {
		return n * 2, true, nil
	})

	var results []int
	err := Drain[int](doubled, func(n int) error {
		results = append(results, n)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 5 {
		t.Fatalf("got %d results, want 5", len(results))
	}
	for i, want := range []int{2, 4, 6, 8, 10} {
		if results[i] != want {
			t.Errorf("results[%d] = %d, want %d", i, results[i], want)
		}
	}
}

func TestMapFilter_Skip(t *testing.T) {
	src := newSliceSource([]int{1, 2, 3, 4, 5})
	evens := NewMapFilter(src, func(n int) (int, bool, error) {
		return n, n%2 == 0, nil
	})

	var results []int
	Drain[int](evens, func(n int) error {
		results = append(results, n)
		return nil
	})
	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}
	if results[0] != 2 || results[1] != 4 {
		t.Errorf("got %v, want [2 4]", results)
	}
}

func TestMapFilter_Error(t *testing.T) {
	src := newSliceSource([]int{1, 2, 3})
	failing := NewMapFilter(src, func(n int) (int, bool, error) {
		if n == 2 {
			return 0, false, fmt.Errorf("boom")
		}
		return n, true, nil
	})

	var results []int
	err := Drain[int](failing, func(n int) error {
		results = append(results, n)
		return nil
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if len(results) != 1 {
		t.Errorf("got %d results before error, want 1", len(results))
	}
}

func TestLimitSource(t *testing.T) {
	src := newSliceSource([]int{1, 2, 3, 4, 5})
	limited := NewLimitSource[int](src, 3)

	var results []int
	Drain[int](limited, func(n int) error {
		results = append(results, n)
		return nil
	})
	if len(results) != 3 {
		t.Fatalf("got %d results, want 3", len(results))
	}
}

func TestChainedPipeline(t *testing.T) {
	src := newSliceSource([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	// limit -> filter evens -> double
	limited := NewLimitSource[int](src, 8)
	evens := NewMapFilter(limited, func(n int) (int, bool, error) {
		return n, n%2 == 0, nil
	})
	doubled := NewMapFilter(evens, func(n int) (int, bool, error) {
		return n * 2, true, nil
	})

	var results []int
	Drain[int](doubled, func(n int) error {
		results = append(results, n)
		return nil
	})

	want := []int{4, 8, 12, 16}
	if len(results) != len(want) {
		t.Fatalf("got %v, want %v", results, want)
	}
	for i := range want {
		if results[i] != want[i] {
			t.Errorf("results[%d] = %d, want %d", i, results[i], want[i])
		}
	}
}

func TestMapFilter_TypeConversion(t *testing.T) {
	src := newSliceSource([]int{1, 2, 3})
	strs := NewMapFilter(src, func(n int) (string, bool, error) {
		return fmt.Sprintf("item_%d", n), true, nil
	})

	var results []string
	Drain[string](strs, func(s string) error {
		results = append(results, s)
		return nil
	})
	if len(results) != 3 {
		t.Fatalf("got %d, want 3", len(results))
	}
	if results[0] != "item_1" {
		t.Errorf("got %q, want item_1", results[0])
	}
}

func TestDrain_EmptySource(t *testing.T) {
	src := newSliceSource([]int{})
	count := 0
	err := Drain[int](src, func(n int) error {
		count++
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Errorf("got %d, want 0", count)
	}
}
