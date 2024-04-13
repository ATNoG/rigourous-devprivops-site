package util

import "strings"

func Map[T1 any, T2 any](arr []T1, mapper func(T1) T2) []T2 {
	new := []T2{}

	for _, e := range arr {
		new = append(new, mapper(e))
	}

	return new
}

func Filter[T any](arr []T, filterFn func(T) bool) []T {
	new := []T{}

	for _, e := range arr {
		if filterFn(e) {
			new = append(new, e)
		}
	}

	return new
}

func Sum(arr []int) int {
	sum := 0

	for _, e := range arr {
		sum += e
	}

	return sum
}

func Btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

func ToHTMLID(input string) string {
	replaced := strings.ReplaceAll(input, " ", "-")
	lowercase := strings.ToLower(replaced)
	return lowercase
}

func Contains[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
