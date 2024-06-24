package polyfill

import (
	"errors"
	"fmt"
)

var ErrRange = errors.New("range error")

type Iterator[T any] interface {
	Next() (T, bool)
}

type iterator[T any] struct {
	index    int
	size     int
	elements []T
	err      error
}

func build[T any](arr []T) *iterator[T] {
	return &iterator[T]{
		size:     len(arr),
		elements: arr,
	}
}

func (iter *iterator[T]) Next() (T, bool) {
	var v T
	if iter.index >= iter.size {
		return v, false
	}
	v = iter.elements[iter.index]
	iter.index++
	return v, true
}

func (iter *iterator[T]) Drop(limit int) *iterator[T] {
	if limit < 0 {
		iter.err = fmt.Errorf("%w: drop limit must be positive", ErrRange)
		return iter
	}

	if iter.err != nil {
		return iter
	}

	iter.index += limit
	return iter
}

func (iter *iterator[T]) Every(fn func(int, T) bool) bool {
	var every bool
	for i := iter.index; i < iter.size; i++ {
		if every = fn(i, iter.elements[i]); !every {
			break
		}
	}
	iter.index = iter.size
	return every
}

func (iter *iterator[T]) Filter(fn func(int, T) bool) *iterator[T] {
	filters := make([]T, 0, iter.size)
	for i := iter.index; i < iter.size; i++ {
		if fn(i, iter.elements[i]) {
			filters = append(filters, iter.elements[i])
		}
	}
	return build(filters)
}

func (iter *iterator[T]) Find(fn func(int, T) bool) (T, bool) {
	var v T
	for i := iter.index; i < iter.size; i++ {
		if fn(i, iter.elements[i]) {
			v = iter.elements[i]
			return v, true
		}
	}
	return v, false
}

func (iter iterator[T]) FlatMap() {
}

func (iter iterator[T]) ForEach(fn func(int, T)) {
	for i := range iter.elements {
		fn(i, iter.elements[i])
	}
}

func (iter *iterator[T]) Map(fn func(int, T) T) *iterator[T] {
	maps := make([]T, 0, iter.size)
	for i := iter.index; i < iter.size; i++ {
		maps = append(maps, fn(i, iter.elements[i]))
	}
	return build(maps)
}

func (iter *iterator[T]) Reduce(fn func(accumulator, current T) T, initial T) T {
	for i := iter.index; i < iter.size; i++ {
		initial = fn(initial, iter.elements[i])
	}
	return initial
}

func (iter *iterator[T]) Some(fn func(int, T) bool) bool {
	for i := iter.index; i < iter.size; i++ {
		if fn(i, iter.elements[i]) {
			return true
		}
	}
	return false
}

func (iter *iterator[T]) Take(limit int) *iterator[T] {
	ans := iter.elements[iter.index : iter.index+limit]
	return build(ans)
}

func (iter *iterator[T]) ToSlice() []T {
	ans := make([]T, iter.size-iter.index)
	copy(ans, iter.elements[iter.index:])
	return ans
}
