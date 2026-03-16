// Package pipeline provides a generic, pull-based filter pipeline.
//
// Values produced by Source.Next() may reference shared buffers and are
// only valid until the next call to Next(). Consumers that need to retain
// a value beyond that must call Copy() on types that implement Copier.
package pipeline

// Source is a pull-based iterator that produces values on demand.
// Next returns the next value, true if valid, and an error.
// When exhausted: (zero, false, nil). On error: (zero, false, err).
type Source[T any] interface {
	Next() (T, bool, error)
}

// Copier is implemented by types whose values reference shared buffers
// and need an explicit deep copy for extended lifetime.
type Copier[T any] interface {
	Copy() T
}

// MapFilter transforms a Source[I] into a Source[O].
// The transform function returns (value, keep, error).
// If keep is false, the element is skipped and the next one is pulled.
type MapFilter[I, O any] struct {
	upstream Source[I]
	fn       func(I) (O, bool, error)
}

func NewMapFilter[I, O any](upstream Source[I], fn func(I) (O, bool, error)) *MapFilter[I, O] {
	return &MapFilter[I, O]{upstream: upstream, fn: fn}
}

func (mf *MapFilter[I, O]) Next() (O, bool, error) {
	var zero O
	for {
		val, ok, err := mf.upstream.Next()
		if !ok || err != nil {
			return zero, false, err
		}
		out, keep, err := mf.fn(val)
		if err != nil {
			return zero, false, err
		}
		if keep {
			return out, true, nil
		}
	}
}

// LimitSource wraps a Source and stops after limit elements.
type LimitSource[T any] struct {
	upstream Source[T]
	limit    int
	count    int
}

func NewLimitSource[T any](upstream Source[T], limit int) *LimitSource[T] {
	return &LimitSource[T]{upstream: upstream, limit: limit}
}

func (ls *LimitSource[T]) Next() (T, bool, error) {
	var zero T
	if ls.count >= ls.limit {
		return zero, false, nil
	}
	val, ok, err := ls.upstream.Next()
	if ok {
		ls.count++
	}
	return val, ok, err
}

// Drain exhausts a Source, calling fn on each element.
func Drain[T any](src Source[T], fn func(T) error) error {
	for {
		val, ok, err := src.Next()
		if !ok {
			return err
		}
		if err != nil {
			return err
		}
		if err := fn(val); err != nil {
			return err
		}
	}
}
