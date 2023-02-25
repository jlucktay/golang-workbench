package main

import (
	"fmt"
	"sort"

	"golang.org/x/exp/constraints"
)

// Map turns a []T1 to a []T2 using a mapping function.
// This function has two type parameters, T1 and T2.
// This works with slices of any type.
func Map[T1, T2 any](s []T1, f func(T1) T2) []T2 {
	r := make([]T2, len(s))

	for i, v := range s {
		r[i] = f(v)
	}

	return r
}

// Reduce reduces a []T1 to a single value using a reduction function.
func Reduce[T1, T2 any](s []T1, initializer T2, f func(T2, T1) T2) T2 {
	r := initializer

	for _, v := range s {
		r = f(r, v)
	}

	return r
}

// Filter filters values from a slice using a filter function.
// It returns a new slice with only the elements of s
// for which f returned true.
func Filter[T any](s []T, f func(T) bool) []T {
	var r []T

	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}

	return r
}

// Merge - receives slices of type T and merges them into a single slice of type T.
func Merge[T any](slices ...[]T) (mergedSlice []T) {
	for _, slice := range slices {
		mergedSlice = append(mergedSlice, slice...)
	}

	return mergedSlice
}

// Includes - given a slice of type T and a value of type T,
// determines whether the value is contained by the slice.
func Includes[T comparable](slice []T, value T) bool {
	for _, el := range slice {
		if el == value {
			return true
		}
	}

	return false
}

// // Sort - sorts given a slice of any orderable type T
// The constraints.Ordered constraint in the Sort() function guarantees that
// the function can sort values of any type supporting the operators <, <=, >=, >.
func Sort[T constraints.Ordered](s []T) {
	sort.Slice(s, func(i, j int) bool {
		return s[i] < s[j]
	})
}

// Keys returns the keys of the map m in a slice.
// The keys will be returned in an unpredictable order.
// This function has two type parameters, K and V.
// Map keys must be comparable, so key has the predeclared
// constraint comparable. Map values can be any type.
func Keys[K comparable, V any](m map[K]V) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}

// Sum sums the values of map containing numeric or float values
func Sum[K comparable, V constraints.Float | constraints.Integer](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}

func main() {
	s := []int{1, 2, 3, 7, 5, 22, 18}
	j := []int{4, 5, 6}

	floats := Map(s, func(i int) float64 { return float64(i) })
	fmt.Println(floats)

	sum := Reduce(s, 0, func(i, j int) int { return i + j })
	fmt.Println(sum)

	evens := Filter(s, func(i int) bool { return i%2 == 0 })
	fmt.Println(evens)

	merged := Merge(s, j)
	fmt.Println(merged)

	i := Includes(s, 22)
	fmt.Println(i)

	Sort(s)
	fmt.Println(s)

	k := Merge(s, j)
	Sort(k)
	odds := Filter(k, func(i int) bool { return i%2 != 0 })
	fmt.Println(odds)

	l := Keys(map[int]int{1: 2, 2: 4})
	fmt.Println(l)

	t := Sum(map[int]int{1: 2, 2: 4})
	fmt.Println(t)
}
