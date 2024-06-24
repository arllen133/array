package polyfill

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

func At[T any](arr []T, index int) (T, bool) {
	var (
		noop T
		size = len(arr)
	)
	if index >= size {
		return noop, false
	}
	if index < 0 {
		index += size
		if index < 0 {
			return noop, false
		}
	}
	return arr[index], true
}

func With[T any](arr []T, index int, element T) []T {
	size := len(arr)
	if index < 0 {
		index += size
		if index < 0 {
			panic(ErrRange)
		}
	} else if index >= size {
		panic(ErrRange)
	}
	ans := make([]T, size)
	copy(ans, arr)
	ans[index] = element
	return ans
}

func Concat[T any](arr []T, src ...[]T) []T {
	size := len(arr)
	for _, src := range src {
		size += len(src)
	}
	ans := make([]T, size)
	copy(ans, arr)
	size = len(arr)
	for _, item := range src {
		copy(ans[size:], item)
		size += len(src)
	}
	return ans
}

func CopyWithin[T any](arr []T, target int, args ...int) []T {
	size := len(arr)
	target = normalizeIndex(target, size)
	if target >= size {
		return arr
	}
	start, end := normalize(size, args...)
	for i := start; i < end; i++ {
		arr[target] = arr[i]
		target++
	}

	return arr
}

func Entries[T any](arr []T) Iterator[T] {
	return build(arr)
}

// Every returns true if all elements in an array pass the test implemented by the provided function
func Every[T any](arr []T, fn func(int, T) bool) bool {
	for i := range arr {
		if !fn(i, arr[i]) {
			return false
		}
	}
	return true
}

// Fill fills an array with a static value
func Fill[T any](arr []T, element T, args ...int) []T {
	size := len(arr)
	start, end := normalize(size, args...)

	for i := start; i < end; i++ {
		arr[i] = element
	}

	return arr
}

// Filter returns a new array with all elements that pass the test implemented by the provided function
func Filter[T any](arr []T, fn func(int, T) bool) []T {
	ans := make([]T, 0, len(arr))
	for i := range arr {
		if fn(i, arr[i]) {
			ans = append(ans, arr[i])
		}
	}
	return ans
}

// Find returns the first element in an array
func Find[T any](arr []T, fn func(int, T) bool) (T, bool) {
	i := FindIndex(arr, fn)
	if i == -1 {
		var noop T
		return noop, false
	}
	return arr[i], true
}

// FindIndex returns the index of an element in an array
func FindIndex[T any](arr []T, fn func(int, T) bool) int {
	for i := range arr {
		if fn(i, arr[i]) {
			return i
		}
	}
	return -1
}

// FindLast returns the last element in an array
func FindLast[T any](arr []T, fn func(int, T) bool) (T, bool) {
	i := FindLastIndex(arr, fn)
	if i == -1 {
		var noop T
		return noop, false
	}
	return arr[i], true
}

// FindLastIndex returns the last index of an element in an array
func FindLastIndex[T any](arr []T, fn func(int, T) bool) int {
	for i := len(arr) - 1; i >= 0; i-- {
		if fn(i, arr[i]) {
			return i
		}
	}
	return -1
}

// Flat returns a new flattened array
func Flat[T any](arr []T, args ...int) []T {
	depth := 1
	if len(args) > 0 {
		depth = args[0]
		if depth < 0 {
			depth = math.MaxInt
		}
	}

	return flattenRecursive(arr, depth)
}

// FlatMap executes a provided function once for each array element and returns a new flattened array
func FlatMap[T any](arr []T, fn func(int, T) T) []T {
	ans := make([]T, 0, len(arr))
	for i := range arr {
		v := fn(i, arr[i])
		switch value := any(v).(type) {
		case []T:
			ans = append(ans, value...)
		default:
			ans = append(ans, v)
		}
	}
	return ans
}

// ForEach executes a provided function once for each array element
func ForEach[T any](arr []T, fn func(int, T)) {
	for i := range arr {
		fn(i, arr[i])
	}
}

func Map[T, V any](arr []T, fn func(int, T) V) []V {
	ans := make([]V, len(arr))
	for i := range arr {
		ans[i] = fn(i, arr[i])
	}
	return ans
}

// Keys returns an array of keys
func Keys[T any](arr []T) []int {
	keys := make([]int, len(arr))
	for i := range arr {
		keys[i] = i
	}
	return keys
}

// Pop removes the last element from an array and returns it
func Pop[T any](arr *[]T) (T, bool) {
	n := len(*arr)
	if n == 0 {
		var noop T
		return noop, false
	}

	v := (*arr)[n-1]
	*arr = (*arr)[:n-1]

	return v, true
}

// Push adds one or more elements to the end of an array and
// returns the new length of the array
func Push[T any](arr *[]T, elements ...T) int {
	*arr = append(*arr, elements...)
	return len(*arr)
}

// Reduce  executes a user-supplied "reducer" callback function on each element of the array,
// in order, passing in the return value from the calculation on the preceding element.
// The final result of running the reducer across all elements of the array is a single value.
func Reduce[T any](arr []T, fn func(accumulator, current T) T, initial T) T {
	for i := range arr {
		initial = fn(initial, arr[i])
	}
	return initial
}

// ReduceRight is the same as Reduce, but reduces the array from right to left.
func ReduceRight[T any](arr []T, fn func(accumulator, current T) T, initial T) T {
	for i := len(arr) - 1; i >= 0; i-- {
		initial = fn(initial, arr[i])
	}
	return initial
}

// Reverse reverses an array
func Reverse[T any](arr []T) []T {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

// Shift removes the first element from an array and returns the removed
// element
func Shift[T any](arr *[]T) (T, bool) {
	n := len(*arr)
	if n == 0 {
		var noop T
		return noop, false
	}
	v := (*arr)[0]
	*arr = (*arr)[1:]
	return v, true
}

// Unshift adds one or more elements to the beginning of an array and
// returns the new length of the array
func Unshift[T any](arr *[]T, elements ...T) int {
	*arr = append(elements, *arr...)
	return len(*arr)
}

// Slice returns a copy of the specified portion of the array
func Slice[T any](arr []T, args ...int) []T {
	size := len(arr)
	switch len(args) {
	case 0:
		return arr
	case 1:
		start := normalizeIndex(args[0], size)
		return arr[start:]
	default:
		start := normalizeIndex(args[0], size)
		end := normalizeIndex(args[1], size)
		if start >= end {
			return []T{}
		}
		return arr[start:end]
	}
}

// Some
func Some[T any](arr []T, fn func(int, T) bool) bool {
	for i := range arr {
		if fn(i, arr[i]) {
			return true
		}
	}
	return false
}

func Sort[T any](arr []T, fn func(i, j int) bool) []T {
	sort.Slice(arr, func(i, j int) bool {
		return fn(i, j)
	})
	return arr
}

func ToSorted[T any](arr []T, fn func(T, T) bool) []T {
	ans := make([]T, len(arr))
	copy(ans, arr)
	sort.Slice(ans, func(i, j int) bool {
		return fn(ans[i], ans[j])
	})
	return ans
}

func ToReserved[T any](arr []T) []T {
	n := len(arr)
	ans := make([]T, n)
	for i, j := n-1, 0; i >= 0; i, j = i-1, j+1 {
		ans[j] = arr[i]
	}
	return ans
}

func Values[T any](arr []T) Iterator[T] {
	return build(arr)
}

func Splice[T any](arr *[]T, start, deleteCount int, items ...T) []T {
	size := len(*arr)
	start = normalizeIndex(start, size)

	if deleteCount < 0 {
		deleteCount = 0
	}
	end := start + deleteCount
	if end > size {
		end = size
	}

	// normalize delete count
	deleteCount = end - start
	// remove element
	var removed []T
	if deleteCount > 0 {
		removed = make([]T, deleteCount)
		copy(removed, (*arr)[start:end])
	} else {
		removed = []T{}
	}

	// add element
	if len(removed) > 0 || len(items) > 0 {
		ans := make([]T, 0, size-deleteCount+len(items))
		ans = append(ans, (*arr)[:start]...)
		ans = append(ans, items...)
		ans = append(ans, (*arr)[end:]...)
		*arr = ans
	}

	return removed
}

////////////////////////////////////////////////////////////////////////////////////////////

func Includes[T comparable](arr []T, element T) bool {
	for i := range arr {
		if arr[i] == element {
			return true
		}
	}
	return false
}

func IndexOf[T comparable](arr []T, element T) int {
	for i := range arr {
		if arr[i] == element {
			return i
		}
	}
	return -1
}

func LastIndexOf[T comparable](arr []T, element T) int {
	for i := len(arr) - 1; i >= 0; i-- {
		if arr[i] == element {
			return i
		}
	}
	return -1
}

////////////////////////////////////////////////////////////////////////////////////////////

type Stringer interface {
	String() string
}

func Join[T Stringer](arr []T, sep string) string {
	switch n := len(arr); n {
	case 0:
		return ""
	case 1:
		return arr[0].String()
	default:
		var sb strings.Builder
		sb.WriteString(arr[0].String())
		for i := 1; i < n; i++ {
			sb.WriteString(sep)
			sb.WriteString(arr[i].String())
		}
		return sb.String()
	}
}

func ToString[T any](arr []T) string {
	switch n := len(arr); n {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf("%v", arr[0])
	default:
		var sb strings.Builder
		for i := 0; i < n; i++ {
			sb.WriteString(fmt.Sprintf("%v", arr[i]))
		}
		return sb.String()
	}
}
