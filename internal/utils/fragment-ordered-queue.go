package utils

import (
	"sort"
)

type OrderedFragment[T any] struct {
	Index   int
	Payload T
}

type FragmentOrderedQueue[T any] struct {
	itemsNumber int
	current     int
	items       []OrderedFragment[T]
}

func CreateFragmentOrderedQueue[T any](itemsNumber int) *FragmentOrderedQueue[T] {
	return &FragmentOrderedQueue[T]{
		itemsNumber: itemsNumber,
		current:     0,
		items:       []OrderedFragment[T]{},
	}
}

func (slice *FragmentOrderedQueue[T]) Enqueue(item OrderedFragment[T]) {
	slice.items = append(slice.items, item)

	sort.Slice(slice.items, func(i, j int) bool {
		return slice.items[i].Index < slice.items[j].Index
	})
}

func (slice *FragmentOrderedQueue[T]) Dequeue() ([]OrderedFragment[T], bool) {
	if len(slice.items) == 0 {
		return []OrderedFragment[T]{}, false
	}

	if slice.items[0].Index != slice.current {
		return []OrderedFragment[T]{}, false
	}

	cutIndex := 1
	for i := cutIndex; i < len(slice.items); i++ {
		if slice.items[i].Index != slice.items[i-1].Index+1 {
			break
		}

		cutIndex += 1
	}

	dequeueItems := slice.items[:cutIndex]
	lastItem := dequeueItems[len(dequeueItems)-1]
	slice.current = lastItem.Index + 1

	slice.items = slice.items[cutIndex:]

	return dequeueItems, lastItem.Index == (slice.itemsNumber - 1)
}
