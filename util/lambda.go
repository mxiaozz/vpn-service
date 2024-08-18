package util

import (
	"slices"
)

type List[T any] []T

func NewList[T any](list []T) List[T] {
	return list
}

func (list List[T]) ForEach(action func(t T)) {
	for _, v := range list {
		action(v)
	}
}

func (list List[T]) Filter(predicate func(t T) bool) List[T] {
	newList := make(List[T], 0)
	for _, v := range list {
		if predicate(v) {
			newList = append(newList, v)
		}
	}
	return newList
}

func (list List[T]) Distinct(genKey func(t T) any) List[T] {
	m := make(map[any]T)
	for _, v := range list {
		key := genKey(v)
		if _, ok := m[key]; !ok {
			m[key] = v
		}
	}

	newList := make(List[T], 0)
	for _, v := range m {
		newList = append(newList, v)
	}
	return newList
}

func (list List[T]) Count() int {
	return len(list)
}

func (list List[T]) Max(cmp func(a, b T) int) T {
	return slices.MaxFunc(list, cmp)
}

func (list List[T]) Min(cmp func(a, b T) int) T {
	return slices.MinFunc(list, cmp)
}

func (list List[T]) Order(cmp func(a, b T) int) List[T] {
	slices.SortFunc(list, cmp)
	return list
}

func (list List[T]) MapToInt(mapper func(t T) int) List[int] {
	return Convert(list, mapper)
}

func (list List[T]) MapToInt64(mapper func(t T) int64) List[int64] {
	return Convert(list, mapper)
}

func (list List[T]) MapToStr(mapper func(t T) string) List[string] {
	return Convert(list, mapper)
}

func Convert[T any, R any](list List[T], mapper func(t T) R) List[R] {
	newList := make(List[R], 0)
	for _, v := range list {
		newList = append(newList, mapper(v))
	}
	return newList
}
